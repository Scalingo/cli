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
	"strconv"
	"strings"
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

	res, socket, err := connectToRunServer(runStruct["attach"].(string))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Fail to attach: %s", res.Status)
	}

	if err := term.MakeRaw(os.Stdin); err != nil {
		return err
	}

	go io.Copy(socket, os.Stdin)
	io.Copy(os.Stdout, socket)

	if err := term.Restore(os.Stdin); err != nil {
		return err
	}

	return nil
}

func connectToRunServer(rawUrl string) (*http.Response, net.Conn, error) {
	params := map[string]string{
		"user_email": auth.Config.Email,
		"user_token": auth.Config.AuthToken,
	}
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return nil, nil, err
	}
	paramsReader := bytes.NewReader(paramsJson)

	req, err := http.NewRequest("POST", rawUrl, paramsReader)
	if err != nil {
		return nil, nil, err
	}

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
