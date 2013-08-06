package cmd

import (
	"fmt"
	"os"
)

func errorQuit(err error) {
	fmt.Printf("[Error] %s\n", err)
	os.Exit(1)
}
