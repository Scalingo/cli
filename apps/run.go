package apps

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	stdio "io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps/run"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/signals"
	"github.com/Scalingo/cli/term"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-scalingo/v6/debug"
)

type RunOpts struct {
	App            string
	DisplayCmd     string
	Silent         bool
	Detached       bool
	Size           string
	Type           string
	Cmd            []string
	CmdEnv         []string
	Files          []string
	StdinCopyFunc  func(stdio.Writer, stdio.Reader) (int64, error)
	StdoutCopyFunc func(stdio.Writer, stdio.Reader) (int64, error)
}

type runContext struct {
	app                     string
	attachURL               string
	scalingoClient          *scalingo.Client
	waitingTextOutputWriter stdio.Writer
	stdinCopyFunc           func(stdio.Writer, stdio.Reader) (int64, error)
	stdoutCopyFunc          func(stdio.Writer, stdio.Reader) (int64, error)
}

func Run(ctx context.Context, opts RunOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	firstReadDone := make(chan struct{})
	runCtx := &runContext{
		app:                     opts.App,
		waitingTextOutputWriter: os.Stderr,
		stdinCopyFunc:           stdio.Copy,
		stdoutCopyFunc:          io.CopyWithFirstReadChan(firstReadDone),
		scalingoClient:          c,
	}
	if opts.Type != "" {
		processes, err := c.AppsContainerTypes(ctx, opts.App)
		if err != nil {
			return errgo.Mask(err)
		}
		for _, p := range processes {
			if p.Name == opts.Type {
				opts.Cmd = strings.Split(p.Command, " ")
			}
		}
		if strings.Join(opts.Cmd, "") == "" {
			return errgo.New("no such type")
		}
	}

	if opts.Size == "" {
		opts.Size = "M"
	}

	if opts.CmdEnv == nil {
		opts.CmdEnv = []string{}
	}
	if opts.Files == nil {
		opts.Files = []string{}
	}
	if opts.Silent {
		runCtx.waitingTextOutputWriter = new(bytes.Buffer)
	}
	if opts.StdinCopyFunc != nil {
		runCtx.stdinCopyFunc = opts.StdinCopyFunc
	}
	if opts.StdoutCopyFunc != nil {
		runCtx.stdoutCopyFunc = opts.StdoutCopyFunc
	}

	env, err := runCtx.buildEnv(opts.CmdEnv)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	err = runCtx.validateFiles(opts.Files)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	runRes, err := c.Run(
		ctx,
		scalingo.RunOpts{
			App:      opts.App,
			Command:  opts.Cmd,
			Env:      env,
			Size:     opts.Size,
			Detached: opts.Detached,
		})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	debug.Printf("%+v\n", runRes)

	if opts.Detached {
		fmt.Printf(
			"Starting one-off '%s' for app '%v'.\n"+
				"Run `scalingo --region %v --app %v logs --filter %v` to get the output\n",
			io.Bold(strings.Join(opts.Cmd, " ")), io.Bold(opts.App),
			config.C.ScalingoRegion, opts.App, runRes.Container.Label,
		)
		return nil
	}

	runCtx.attachURL = runRes.AttachURL
	debug.Println("Run Service URL is", runCtx.attachURL)

	if len(opts.Files) > 0 {
		err := runCtx.uploadFiles(ctx, runCtx.attachURL+"/files", opts.Files)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(
		runCtx.waitingTextOutputWriter,
		"-----> Connecting to container [%v]...  ",
		runRes.Container.Label,
	)

	attachSpinner := io.NewSpinner(runCtx.waitingTextOutputWriter)
	attachSpinner.PostHook = func() {
		var displayCmd string
		if opts.DisplayCmd != "" {
			displayCmd = opts.DisplayCmd
		} else {
			displayCmd = strings.Join(opts.Cmd, " ")
		}
		fmt.Fprintf(runCtx.waitingTextOutputWriter, "\n-----> Process '%v' is starting...  ", displayCmd)
	}
	go attachSpinner.Start()

	res, socket, err := runCtx.connectToRunServer(ctx)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	if res.StatusCode != http.StatusOK {
		return errgo.Newf("Fail to attach: %s", res.Status)
	}

	if term.IsATTY(os.Stdin) {
		if err := term.MakeRaw(os.Stdin); err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}

	stopSignalsMonitoring := make(chan bool)
	defer close(stopSignalsMonitoring)

	signals.CatchQuitSignals = false
	signals := run.NotifiedSignals()

	go func() {
		defer close(signals)

		for {
			select {
			case s := <-signals:
				run.HandleSignal(ctx, runCtx.scalingoClient, s, socket, runCtx.attachURL)
			case <-stopSignalsMonitoring:
				signal.Stop(signals)
				return
			}
		}
	}()

	attachSpinner.Stop()
	startSpinner := io.NewSpinnerWithStopChan(runCtx.waitingTextOutputWriter, firstReadDone)
	// This method will be executed after first read
	startSpinner.PostHook = func() {
		go run.NotifyTermSizeUpdate(signals)
		fmt.Fprintf(runCtx.waitingTextOutputWriter, "\n\n")
	}
	go startSpinner.Start()

	go func() {
		_, err := runCtx.stdinCopyFunc(socket, os.Stdin)
		if err != nil {
			debug.Println("error after reading stdin", err)
		} else {
			// Send EOT when stdin returns
			// 'scalingo run < file'
			socket.Write([]byte("\x04"))
		}
	}()

	_, err = runCtx.stdoutCopyFunc(os.Stdout, socket)

	stopSignalsMonitoring <- true

	if term.IsATTY(os.Stdin) {
		if err := term.Restore(os.Stdin); err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}

	exitCode, err := runCtx.exitCode(ctx)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	os.Exit(exitCode)
	return nil
}

func (ctx *runContext) buildEnv(cmdEnv []string) (map[string]string, error) {
	env := map[string]string{
		"TERM":      os.Getenv("TERM"),
		"CLIENT_OS": runtime.GOOS,
	}

	for _, cmdVar := range cmdEnv {
		v := strings.SplitN(cmdVar, "=", 2)
		if len(v) != 2 || len(v[0]) == 0 || len(v[1]) == 0 {
			return nil, fmt.Errorf("Invalid environment, format is '--env VARIABLE=value'")
		}
		env[v[0]] = v[1]
	}
	return env, nil
}

func (runCtx *runContext) exitCode(ctx context.Context) (int, error) {
	if runCtx.attachURL == "" {
		return -1, errgo.New("No attach URL to connect to")
	}

	req, err := http.NewRequest("GET", runCtx.attachURL+"/wait", nil)
	if err != nil {
		return -1, errgo.Mask(err, errgo.Any)
	}

	token, err := runCtx.scalingoClient.GetAccessToken(ctx)
	if err != nil {
		return -1, errgo.Notef(err, "fail to generate auth")
	}

	req.SetBasicAuth("", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, errgo.Mask(err)
	}
	defer res.Body.Close()

	body, err := stdio.ReadAll(res.Body)
	if err != nil {
		return -1, errgo.Notef(err, "fail to read body when getting exit code")
	}
	debug.Println("exit code body:", string(body))

	if res.StatusCode == http.StatusRequestTimeout {
		fmt.Println()
		io.Warning("Connection timed out due to inactivity, one-off aborted.")
		io.Info("Data should be sent to/from the container regularly to avoid such timeout")
		fmt.Println()
		io.Info("If you need to run long background tasks, the '--detached' should be used")
		io.Info("In this case, output will be available in the main logs of your application:")
		io.Info(io.Gray(io.Bold(fmt.Sprintf("  $ scalingo -a %s logs", runCtx.app))))
		return -127, nil
	}

	waitRes := map[string]int{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&waitRes)
	if err != nil {
		return -1, errgo.Notef(err, "invalid response when getting exit code")
	}

	return waitRes["exit_code"], nil
}

func (runCtx *runContext) connectToRunServer(ctx context.Context) (*http.Response, net.Conn, error) {
	if runCtx.attachURL == "" {
		return nil, nil, errgo.New("No attach URL to connect to")
	}

	req, err := http.NewRequest("CONNECT", runCtx.attachURL, nil)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}
	token, err := runCtx.scalingoClient.GetAccessToken(ctx)

	if err != nil {
		return nil, nil, errgo.Notef(err, "fail to generate auth")
	}
	req.SetBasicAuth("", token)

	url, err := url.Parse(runCtx.attachURL)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	dial, err := net.Dial("tcp", url.Host)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	var conn *httputil.ClientConn
	if url.Scheme == "https" {
		host := strings.Split(url.Host, ":")[0]
		tlsConfig := config.TlsConfig.Clone()
		tlsConfig.ServerName = host
		tlsConn := tls.Client(dial, tlsConfig)
		conn = httputil.NewClientConn(tlsConn, nil)
	} else if url.Scheme == "http" {
		conn = httputil.NewClientConn(dial, nil)
	} else {
		return nil, nil, errgo.Newf("Invalid scheme format %s", url.Scheme)
	}

	res, err := conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		if err, ok := err.(*net.OpError); ok {
			if err.Err.Error() == "record overflow" {
				return nil, nil, errgo.Newf(
					"Fail to create a secure connection to Scalingo server\n"+
						"The encountered error is: %v (ID: CLI-1001)\n"+
						"Your firewall or proxy may block the connection to %s",
					err, url.Host,
				)
			}
		}
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	connection, _ := conn.Hijack()
	return res, connection, nil
}

func (runCtx *runContext) validateFiles(files []string) error {
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			return errgo.Notef(err, "can't upload %s", file)
		}
	}
	return nil
}

func (runCtx *runContext) uploadFiles(ctx context.Context, endpoint string, files []string) error {
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return errgo.Notef(err, "can't stat file %s", file)
		}
		relPath := file
		file, err = filepath.Abs(relPath)
		if err != nil {
			return errgo.Notef(err, "impossible to get absolute path of %s", relPath)
		}
		if stat.IsDir() {
			dir := file
			file, err = runCtx.compressDir(dir)
			if err != nil {
				return errgo.Notef(err, "fail to compress directory %s", dir)
			}
		}
		err = runCtx.uploadFile(ctx, endpoint, file)
		if err != nil {
			return errgo.Notef(err, "fail to upload file %s", file)
		}
	}
	return nil
}

func (runCtx *runContext) compressDir(dir string) (string, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "job-file")
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	fd, err := os.OpenFile(filepath.Join(tmpDir, filepath.Base(dir)+".tar"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	fmt.Fprintln(runCtx.waitingTextOutputWriter, "Compressing directory", dir, "to", fd.Name())

	err = runCtx.createTarArchive(fd, dir)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	file, err := runCtx.compressToGzip(fd.Name())
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	return file, nil
}

func (runCtx *runContext) createTarArchive(fd *os.File, dir string) error {
	tarFd := tar.NewWriter(fd)
	defer tarFd.Close()
	err := filepath.Walk(dir, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		if info.IsDir() {
			return nil
		}
		tarHeader, err := tar.FileInfoHeader(info, name)
		if err != nil {
			return fmt.Errorf("fail to build tar header: %v", err)
		}
		err = tarFd.WriteHeader(tarHeader)
		if err != nil {
			return fmt.Errorf("fail to write tar header: %v", err)
		}
		fileFd, err := os.OpenFile(name, os.O_RDONLY, 0600)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		_, err = stdio.Copy(tarFd, fileFd)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		return nil
	})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}

func (runCtx *runContext) compressToGzip(file string) (string, error) {
	fdSource, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	defer fdSource.Close()
	fdDest, err := os.OpenFile(file+".gz", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	defer fdDest.Close()
	writer := gzip.NewWriter(fdDest)
	defer writer.Close()

	_, err = stdio.Copy(writer, fdSource)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	return fdDest.Name(), nil
}

func (runCtx *runContext) uploadFile(ctx context.Context, endpoint string, file string) error {
	body := new(bytes.Buffer)
	name := filepath.Base(file)
	multipartFile := multipart.NewWriter(body)
	writer, err := multipartFile.CreateFormFile("file", name)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	fd, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	_, err = stdio.Copy(writer, fd)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	err = fd.Close()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	err = multipartFile.Close()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	token, err := runCtx.scalingoClient.GetAccessToken(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to generate token")
	}
	req.SetBasicAuth("", token)

	req.Header.Set("Content-Type", multipartFile.FormDataContentType())

	fmt.Fprintln(runCtx.waitingTextOutputWriter, "Upload", file, "to container.")
	debug.Println("Endpoint:", req.URL)

	res, err := httpclient.Do(req)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		b, _ := stdio.ReadAll(res.Body)
		return errgo.Newf("Invalid return code %v (%s)", res.Status, strings.TrimSpace(string(b)))
	}
	return nil
}
