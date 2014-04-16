package cmd

import (
	"fmt"
	"os"
)

func errorQuit(err error) {
	fmt.Printf("[Error] %v\n", err)
	os.Exit(1)
}
