package run

import "net"
import "os"

func NotifiedSignals() chan os.Signal {
	signals := make(chan os.Signal)
	return signals
}

func NofityTermSizeUpdate(signals chan os.Signal) {
	return
}

func HandleSignal(s os.Signal, socket net.Conn, runUrl string) {
	return
}
