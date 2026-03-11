//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package run

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/term"

	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/debug"
	"github.com/Scalingo/go-utils/errors/v3"
)

type UpdateTtyParams struct {
	Width  string `json:"width"`
	Height string `json:"height"`
}

func NotifiedSignals() chan os.Signal {
	signals := make(chan os.Signal, 1)
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
		return errors.Wrapf(ctx, err, "fail to get the terminal size")
	}

	params := UpdateTtyParams{
		strconv.Itoa(cols),
		strconv.Itoa(lines),
	}
	paramsJSON, _ := json.Marshal(&params)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(paramsJSON))
	if err != nil {
		return errors.Wrap(ctx, err, "build terminal size update request")
	}
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get access token")
	}
	req.SetBasicAuth("", token)
	debug.Printf("Updating TTY Size: PUT %v %+v", url, params)

	res, err := httpclient.Do(req)
	if err != nil {
		return errors.Wrap(ctx, err, "send terminal size update request")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.Newf(ctx, "Invalid error code from run server: %s", res.Status)
	}

	return nil
}
