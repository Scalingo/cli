package signals

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Handle() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for sig := range signals {
		fmt.Printf("%v catched, aborting…\n", sig)
		os.Exit(-127)
	}
}
