package apps

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
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
	"strings"
	"syscall"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/cli/term"
	"gopkg.in/errgo.v1"
)

func Run(app string, command []string, cmdEnv []string, files []string) error {
	env, err := buildEnv(cmdEnv)
	if err != nil {
		return err
	}

	err = validateFiles(files)
	if err != nil {
		return err
	}

	res, err := api.Run(app, command, env)
	if err != nil {
		return err
	}
	runStruct := make(map[string]interface{})
	ReadJson(res.Body, &runStruct)
	debug.Printf("%+v\n", runStruct)

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("application %s not found", app)
	}

	container := runStruct["container"].(map[string]interface{})

	if _, ok := container["attach"]; !ok {
		return fmt.Errorf("Unexpected answer from server")
	}

	runUrl := container["attach"].(string)
	debug.Println("Run Service URL is", runUrl)

	if len(files) > 0 {
		err := uploadFiles(runUrl+"/files", files)
		if err != nil {
			return err
		}
	}

	res, socket, err := connectToRunServer(runUrl)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Fail to attach: %s", res.Status)
	}

	if err := term.MakeRaw(os.Stdin); err != nil {
		return err
	}

	stopSignalsMonitoring := make(chan bool)
	defer close(stopSignalsMonitoring)

	go func() {
		signals := make(chan os.Signal)
		defer close(signals)

		signal.Notify(signals,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGTSTP,
			syscall.SIGWINCH,
		)

		go func() { signals <- syscall.SIGWINCH }()
		for {
			select {
			case s := <-signals:
				switch s {
				case syscall.SIGINT:
					socket.Write([]byte{0x03})
				case syscall.SIGQUIT:
					socket.Write([]byte{0x1c})
				case syscall.SIGTSTP:
					socket.Write([]byte{0x1a})
				case syscall.SIGWINCH:
					err := updateTtySize(runUrl)
					if err != nil {
						fmt.Println("WARN: Error when updating terminal size:", err)
					}
				}
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
		return err
	}

	return nil
}

func buildEnv(cmdEnv []string) (map[string]string, error) {
	env := map[string]string{
		"TERM": os.Getenv("TERM"),
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
	req, err := http.NewRequest("POST", rawUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	req.SetBasicAuth("", api.CurrentUser.AuthToken)

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, nil, err
	}

	dial, err := net.Dial("tcp", url.Host)
	if err != nil {
		return nil, nil, err
	}

	var conn *httputil.ClientConn
	if url.Scheme == "https" {
		host := strings.Split(url.Host, ":")[0]
		tls_conn := tls.Client(dial, config.GenTLSConfig(host))
		conn = httputil.NewClientConn(tls_conn, nil)
	} else if url.Scheme == "http" {
		conn = httputil.NewClientConn(dial, nil)
	} else {
		return nil, nil, fmt.Errorf("Invalid scheme format %s", url.Scheme)
	}

	res, err := conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		if err, ok := err.(*net.OpError); ok {
			if err.Err.Error() == "record overflow" {
				return nil, nil, fmt.Errorf(
					"Fail to create a secure connection to Scalingo server\n"+
						"The encountered error is: %v (ID: CLI-1001)\n"+
						"Your firewall or proxy may block the connection to %s",
					err, url.Host,
				)
			}
		}
		return nil, nil, err
	}

	connection, _ := conn.Hijack()
	return res, connection, nil
}

type UpdateTtyParams struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

func updateTtySize(url string) error {
	cols, err := term.Cols()
	if err != nil {
		return err
	}
	lines, err := term.Lines()
	if err != nil {
		return err
	}

	params := UpdateTtyParams{
		fmt.Sprintf("%d", cols),
		fmt.Sprintf("%d", lines),
	}
	paramsJson, _ := json.Marshal(&params)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(paramsJson))
	if err != nil {
		return err
	}
	req.SetBasicAuth("", api.CurrentUser.AuthToken)

	res, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Invalid error code from run server: %s", res.Status)
	}

	return nil
}

func validateFiles(files []string) error {
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("can't upload %s: %v", file, err)
		}
	}
	return nil
}

func uploadFiles(endpoint string, files []string) error {
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("can't stat file %s: %v", file, err)
		}
		relPath := file
		file, err = filepath.Abs(relPath)
		if err != nil {
			return fmt.Errorf("impossible to get absolute path of %s", relPath)
		}
		if stat.IsDir() {
			dir := file
			file, err = compressDir(dir)
			if err != nil {
				return fmt.Errorf("fail to compress directory %s: %v", dir, err)
			}
		}
		err = uploadFile(endpoint, file)
		if err != nil {
			return fmt.Errorf("fail to upload file %s: %v", file, err)
		}
	}
	return nil
}

func compressDir(dir string) (string, error) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "job-file")
	if err != nil {
		return "", err
	}
	fd, err := os.OpenFile(filepath.Join(tmpDir, filepath.Base(dir)+".tar"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	fmt.Println("Compressing directory", dir, "to", fd.Name())

	err = createTarArchive(fd, dir)
	if err != nil {
		return "", err
	}

	file, err := compressToGzip(fd.Name())
	if err != nil {
		return "", err
	}

	return file, nil
}

func createTarArchive(fd *os.File, dir string) error {
	tarFd := tar.NewWriter(fd)
	defer tarFd.Close()
	err := filepath.Walk(dir, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
			return err
		}
		_, err = io.Copy(tarFd, fileFd)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func compressToGzip(file string) (string, error) {
	fdSource, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}
	defer fdSource.Close()
	fdDest, err := os.OpenFile(file+".gz", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer fdDest.Close()
	writer := gzip.NewWriter(fdDest)
	defer writer.Close()

	_, err = io.Copy(writer, fdSource)
	if err != nil {
		return "", err
	}

	return fdDest.Name(), nil
}

func uploadFile(endpoint string, file string) error {
	body := new(bytes.Buffer)
	name := filepath.Base(file)
	multipartFile := multipart.NewWriter(body)
	writer, err := multipartFile.CreateFormFile("file", name)
	if err != nil {
		return errgo.Mask(err)
	}

	fd, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = io.Copy(writer, fd)
	if err != nil {
		return errgo.Mask(err)
	}

	err = fd.Close()
	if err != nil {
		return err
	}
	err = multipartFile.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return errgo.Mask(err)
	}
	req.SetBasicAuth("", api.CurrentUser.AuthToken)

	req.Header.Set("Content-Type", multipartFile.FormDataContentType())

	fmt.Println("Upload", file, "to container.")
	debug.Println("Endpoint:", req.URL)

	res, err := httpclient.Do(req)
	if err != nil {
		return errgo.Mask(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		return errgo.Newf("Invalid return code %v (%s)", res.Status, strings.TrimSpace(string(b)))
	}
	return nil
}
