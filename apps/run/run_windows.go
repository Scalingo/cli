package run

import (
	"net"
	"os"

	scalingo "github.com/Scalingo/go-scalingo/v4"
)

func NotifiedSignals() chan os.Signal {
	signals := make(chan os.Signal)
	return signals
}

func NotifyTermSizeUpdate(signals chan os.Signal) {
	return
}

func HandleSignal(c *scalingo.Client, s os.Signal, socket net.Conn, runURL string) {
	return
}
