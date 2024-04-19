package run

import (
	"context"
	"net"
	"os"

	scalingo "github.com/Scalingo/go-scalingo/v7"
)

func NotifiedSignals() chan os.Signal {
	signals := make(chan os.Signal)
	return signals
}

func NotifyTermSizeUpdate(signals chan os.Signal) {
	return
}

func HandleSignal(ctx context.Context, c *scalingo.Client, s os.Signal, socket net.Conn, runURL string) {
	return
}
