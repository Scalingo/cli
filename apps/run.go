package apps

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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
	"github.com/Scalingo/go-scalingo/v4"
	"github.com/Scalingo/go-scalingo/v4/debug"
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

func Run(opts RunOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	firstReadDone := make(chan struct{})
	ctx := &runContext{
		app:                     opts.App,
		waitingTextOutputWriter: os.Stderr,
		stdinCopyFunc:           stdio.Copy,
		stdoutCopyFunc:          io.CopyWithFirstReadChan(firstReadDone),
		scalingoClient:          c,
	}
	if opts.Type != "" {
		processes, err := c.AppsContainerTypes(opts.App)
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
		ctx.waitingTextOutputWriter = new(bytes.Buffer)
	}
	if opts.StdinCopyFunc != nil {
		ctx.stdinCopyFunc = opts.StdinCopyFunc
	}
	if opts.StdoutCopyFunc != nil {
		ctx.stdoutCopyFunc = opts.StdoutCopyFunc
	}

	env, err := ctx.buildEnv(opts.CmdEnv)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	err = ctx.validateFiles(opts.Files)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	runRes, err := c.Run(scalingo.RunOpts{
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
			"Starting one-off '%s' for app '%v'.\nRun `scalingo -a %v logs -F %v` to get the output\n",
			io.Bold(strings.Join(opts.Cmd, " ")), io.Bold(opts.App), opts.App, runRes.Container.Label,
		)
		return nil
	}

	ctx.attachURL = runRes.AttachURL
	debug.Println("Run Service URL is", ctx.attachURL)

	if len(opts.Files) > 0 {
		err := ctx.uploadFiles(ctx.attachURL+"/files", opts.Files)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(
		ctx.waitingTextOutputWriter,
		"-----> Connecting to container [%v]...  ",
		runRes.Container.Label,
	)

	attachSpinner := io.NewSpinner(ctx.waitingTextOutputWriter)
	attachSpinner.PostHook = func() {
		var displayCmd string
		if opts.DisplayCmd != "" {
			displayCmd = opts.DisplayCmd
		} else {
			displayCmd = strings.Join(opts.Cmd, " ")
		}
		fmt.Fprintf(ctx.waitingTextOutputWriter, "\n-----> Process '%v' is starting...  ", displayCmd)
	}
	go attachSpinner.Start()

	res, socket, err := ctx.connectToRunServer()
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
				run.HandleSignal(ctx.scalingoClient, s, socket, ctx.attachURL)
			case <-stopSignalsMonitoring:
				signal.Stop(signals)
				return
			}
		}
	}()

	attachSpinner.Stop()
	startSpinner := io.NewSpinnerWithStopChan(ctx.waitingTextOutputWriter, firstReadDone)
	// This method will be executed after first read
	startSpinner.PostHook = func() {
		go run.NotifyTermSizeUpdate(signals)
		fmt.Fprintf(ctx.waitingTextOutputWriter, "\n\n")
	}
	go startSpinner.Start()

	go func() {
		_, err := ctx.stdinCopyFunc(socket, os.Stdin)
		if err != nil {
			debug.Println("error after reading stdin", err)
		} else {
			// Send EOT when stdin returns
			// 'scalingo run < file'
			socket.Write([]byte("\x04"))
		}
	}()

	_, err = ctx.stdoutCopyFunc(os.Stdout, socket)

	stopSignalsMonitoring <- true

	if term.IsATTY(os.Stdin) {
		if err := term.Restore(os.Stdin); err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	}

	exitCode, err := ctx.exitCode()
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

func (ctx *runContext) exitCode() (int, error) {
	if ctx.attachURL == "" {
		return -1, errgo.New("No attach URL to connect to")
	}

	req, err := http.NewRequest("GET", ctx.attachURL+"/wait", nil)
	if err != nil {
		return -1, errgo.Mask(err, errgo.Any)
	}

	token, err := ctx.scalingoClient.GetAccessToken()
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
		io.Info(io.Gray(io.Bold(fmt.Sprintf("  $ scalingo -a %s logs", ctx.app))))
		return -127, nil
	}

	waitRes := map[string]int{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&waitRes)
	if err != nil {
		return -1, errgo.Notef(err, "invalid response when getting exit code")
	}

	return waitRes["exit_code"], nil
}

func (ctx *runContext) connectToRunServer() (*http.Response, net.Conn, error) {
	if ctx.attachURL == "" {
		return nil, nil, errgo.New("No attach URL to connect to")
	}

	req, err := http.NewRequest("CONNECT", ctx.attachURL, nil)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}
	token, err := ctx.scalingoClient.GetAccessToken()

	if err != nil {
		return nil, nil, errgo.Notef(err, "fail to generate auth")
	}
	req.SetBasicAuth("", token)

	url, err := url.Parse(ctx.attachURL)
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

func (ctx *runContext) validateFiles(files []string) error {
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			return errgo.Notef(err, "can't upload %s", file)
		}
	}
	return nil
}

func (ctx *runContext) uploadFiles(endpoint string, files []string) error {
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
			file, err = ctx.compressDir(dir)
			if err != nil {
				return errgo.Notef(err, "fail to compress directory %s", dir)
			}
		}
		err = ctx.uploadFile(endpoint, file)
		if err != nil {
			return errgo.Notef(err, "fail to upload file %s", file)
		}
	}
	return nil
}

func (ctx *runContext) compressDir(dir string) (string, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "job-file")
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	fd, err := os.OpenFile(filepath.Join(tmpDir, filepath.Base(dir)+".tar"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	fmt.Fprintln(ctx.waitingTextOutputWriter, "Compressing directory", dir, "to", fd.Name())

	err = ctx.createTarArchive(fd, dir)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	file, err := ctx.compressToGzip(fd.Name())
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	return file, nil
}

func (ctx *runContext) createTarArchive(fd *os.File, dir string) error {
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

func (ctx *runContext) compressToGzip(file string) (string, error) {
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

func (ctx *runContext) uploadFile(endpoint string, file string) error {
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

	token, err := ctx.scalingoClient.GetAccessToken()
	if err != nil {
		return errgo.Notef(err, "fail to generate token")
	}
	req.SetBasicAuth("", token)

	req.Header.Set("Content-Type", multipartFile.FormDataContentType())

	fmt.Fprintln(ctx.waitingTextOutputWriter, "Upload", file, "to container.")
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
