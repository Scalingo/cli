// +build darwin dragonfly freebsd linux netbsd openbsd

package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/heroku/hk/term"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/go-scalingo/v4"
	"github.com/Scalingo/go-scalingo/v4/debug"
)

type UpdateTtyParams struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

func NotifiedSignals() chan os.Signal {
	signals := make(chan os.Signal)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTSTP,
		syscall.SIGWINCH,
	)
	return signals
}

func NotifyTermSizeUpdate(signals chan os.Signal) {
	signals <- syscall.SIGWINCH
}

func HandleSignal(c *scalingo.Client, s os.Signal, socket net.Conn, runUrl string) {
	switch s {
	case syscall.SIGINT:
		socket.Write([]byte{0x03})
	case syscall.SIGQUIT:
		socket.Write([]byte{0x1c})
	case syscall.SIGTSTP:
		socket.Write([]byte{0x1a})
	case syscall.SIGWINCH:
		err := updateTtySize(c, runUrl)
		if err != nil {
			debug.Println("WARN: Error when updating terminal size:", err)
		}
	}
}

func updateTtySize(c *scalingo.Client, url string) error {
	cols, err := term.Cols()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	lines, err := term.Lines()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	params := UpdateTtyParams{
		fmt.Sprintf("%d", cols),
		fmt.Sprintf("%d", lines),
	}
	paramsJson, err := json.Marshal(&params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(paramsJson))
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	token, err := c.GetAccessToken()
	if err != nil {
		return errgo.Notef(err, "fail to get access token")
	}
	req.SetBasicAuth("", token)
	debug.Printf("Updating TTY Size: PUT %v %+v", url, params)

	res, err := httpclient.Do(req)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errgo.Newf("Invalid error code from run server: %s", res.Status)
	}

	return nil
}
