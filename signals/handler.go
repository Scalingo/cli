package signals

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	CatchQuitSignals = true
)

func Handle() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for sig := range signals {
		if CatchQuitSignals {
			fmt.Printf("%v catched, aborting…\n", sig)
			os.Exit(-127)
		}
	}
}
