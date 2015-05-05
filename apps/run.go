package apps

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
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

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/apps/run"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/cli/signals"
	"github.com/Scalingo/cli/term"
)

func Run(app string, command []string, cmdEnv []string, files []string) error {
	env, err := buildEnv(cmdEnv)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	err = validateFiles(files)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	res, err := api.Run(app, command, env)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	runStruct := make(map[string]interface{})
	ReadJson(res.Body, &runStruct)
	debug.Printf("%+v\n", runStruct)

	if res.StatusCode == http.StatusNotFound {
		return errgo.Newf("application %s not found", app)
	}

	attachURL, ok := runStruct["attach_url"].(string)
	if !ok {
		return errgo.New("unexpected answer from server")
	}

	debug.Println("Run Service URL is", attachURL)

	if len(files) > 0 {
		err := uploadFiles(attachURL+"/files", files)
		if err != nil {
			return err
		}
	}

	res, socket, err := connectToRunServer(attachURL)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	if res.StatusCode != http.StatusOK {
		return errgo.Newf("Fail to attach: %s", res.Status)
	}

	if err := term.MakeRaw(os.Stdin); err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	stopSignalsMonitoring := make(chan bool)
	defer close(stopSignalsMonitoring)

	go func() {
		signals.CatchQuitSignals = false
		signals := run.NotifiedSignals()
		defer close(signals)

		go run.NofityTermSizeUpdate(signals)
		for {
			select {
			case s := <-signals:
				run.HandleSignal(s, socket, attachURL)
			case <-stopSignalsMonitoring:
				signal.Stop(signals)
				return
			}
		}
	}()

	go io.Copy(socket, os.Stdin)
	io.Copy(os.Stdout, socket)

	stopSignalsMonitoring <- true

	if err := term.Restore(os.Stdin); err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	return nil
}

func buildEnv(cmdEnv []string) (map[string]string, error) {
	env := map[string]string{
		"TERM":      os.Getenv("TERM"),
		"CLIENT_OS": runtime.GOOS,
	}

	for _, cmdVar := range cmdEnv {
		v := strings.Split(cmdVar, "=")
		if len(v[0]) == 0 || len(v[1]) == 0 {
			return nil, fmt.Errorf("Invalid environment, format is '--env VARIABLE=value'")
		}
		env[v[0]] = v[1]
	}
	return env, nil
}

func connectToRunServer(rawUrl string) (*http.Response, net.Conn, error) {
	req, err := http.NewRequest("CONNECT", rawUrl, nil)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}
	req.SetBasicAuth("", api.CurrentUser.AuthToken)

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	dial, err := net.Dial("tcp", url.Host)
	if err != nil {
		return nil, nil, errgo.Mask(err, errgo.Any)
	}

	var conn *httputil.ClientConn
	if url.Scheme == "https" {
		tls_conn := tls.Client(dial, config.TlsConfig)
		conn = httputil.NewClientConn(tls_conn, nil)
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

func validateFiles(files []string) error {
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			return errgo.Notef(err, "can't upload %s", file)
		}
	}
	return nil
}

func uploadFiles(endpoint string, files []string) error {
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
			file, err = compressDir(dir)
			if err != nil {
				return errgo.Notef(err, "fail to compress directory %s", dir)
			}
		}
		err = uploadFile(endpoint, file)
		if err != nil {
			return errgo.Notef(err, "fail to upload file %s", file)
		}
	}
	return nil
}

func compressDir(dir string) (string, error) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "job-file")
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	fd, err := os.OpenFile(filepath.Join(tmpDir, filepath.Base(dir)+".tar"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	fmt.Println("Compressing directory", dir, "to", fd.Name())

	err = createTarArchive(fd, dir)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	file, err := compressToGzip(fd.Name())
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	return file, nil
}

func createTarArchive(fd *os.File, dir string) error {
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
		_, err = io.Copy(tarFd, fileFd)
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

func compressToGzip(file string) (string, error) {
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

	_, err = io.Copy(writer, fdSource)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}

	return fdDest.Name(), nil
}

func uploadFile(endpoint string, file string) error {
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

	_, err = io.Copy(writer, fd)
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
	req.SetBasicAuth("", api.CurrentUser.AuthToken)

	req.Header.Set("Content-Type", multipartFile.FormDataContentType())

	fmt.Println("Upload", file, "to container.")
	debug.Println("Endpoint:", req.URL)

	res, err := httpclient.Do(req)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		return errgo.Newf("Invalid return code %v (%s)", res.Status, strings.TrimSpace(string(b)))
	}
	return nil
}
