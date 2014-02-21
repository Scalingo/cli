package apps

import (
	"appsdeck/api"
	"appsdeck/auth"
	"appsdeck/config"
	"appsdeck/term"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func Run(app string, command []string, cmdEnv []string) error {

	cols, err := term.Cols()
	if err != nil {
		log.Fatal(err)
	}
	lines, err := term.Lines()
	if err != nil {
		log.Fatal(err)
	}

	env := map[string]string{
		"COLUMNS": strconv.Itoa(cols),
		"LINES":   strconv.Itoa(lines),
		"TERM":    os.Getenv("TERM"),
	}

	for _, cmdVar := range cmdEnv {
		v := strings.Split(cmdVar, "=")
		if len(v[0]) == 0 || len(v[1]) == 0 {
			return fmt.Errorf("Invalid environment")
		}
		env[v[0]] = v[1]
	}

	res, err := api.Run(app, command, env)
	if err != nil {
		return err
	}
	runStruct := make(map[string]interface{})
	ReadJson(res.Body, &runStruct)

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Not authorized")
	}

	if _, ok := runStruct["attach"]; !ok {
		return fmt.Errorf("Unexpected answer from server")
	}

	runUrl := runStruct["attach"].(string)

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

func connectToRunServer(rawUrl string) (*http.Response, net.Conn, error) {
	req, err := http.NewRequest("POST", rawUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	auth.AddHeaders(req)

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
		tls_conn := tls.Client(dial, config.TlsConfig)
		conn = httputil.NewClientConn(tls_conn, nil)
	} else if url.Scheme == "http" {
		conn = httputil.NewClientConn(dial, nil)
	} else {
		return nil, nil, fmt.Errorf("Invalid scheme format %s", url.Scheme)
	}

	res, err := conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		return nil, nil, err
	}

	connection, _ := conn.Hijack()
	return res, connection, nil
}

type UpdateTtyParams struct {
	Width  string `json: "width"`
	Height string `json: "height"`
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
	auth.AddHeaders(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Invalid error code from run server: %s", res.Status)
	}

	return nil
}
