//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd
// +build darwin dragonfly freebsd linux netbsd openbsd

package run

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-scalingo/v9/debug"
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

func HandleSignal(ctx context.Context, c *scalingo.Client, s os.Signal, socket net.Conn, runURL string) {
	switch s {
	case syscall.SIGINT:
		socket.Write([]byte{0x03})
	case syscall.SIGQUIT:
		socket.Write([]byte{0x1c})
	case syscall.SIGTSTP:
		socket.Write([]byte{0x1a})
	case syscall.SIGWINCH:
		err := updateTtySize(ctx, c, runURL)
		if err != nil {
			debug.Println("WARN: Error when updating terminal size:", err)
		}
	}
}

func updateTtySize(ctx context.Context, c *scalingo.Client, url string) error {
	cols, lines, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return errgo.Notef(err, "fail to get the terminal size")
	}

	params := UpdateTtyParams{
		fmt.Sprintf("%d", cols),
		fmt.Sprintf("%d", lines),
	}
	paramsJSON, _ := json.Marshal(&params)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(paramsJSON))
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	token, err := c.GetAccessToken(ctx)
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
