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

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/heroku/hk/term"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/httpclient"
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

func NofityTermSizeUpdate(signals chan os.Signal) {
	signals <- syscall.SIGWINCH
}

func HandleSignal(s os.Signal, socket net.Conn, runUrl string) {
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
}

func updateTtySize(url string) error {
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
	paramsJson, _ := json.Marshal(&params)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(paramsJson))
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	req.SetBasicAuth("", scalingo.CurrentUser.AuthToken)

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
