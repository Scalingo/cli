package cmd

import (
	"fmt"
	"github.com/Appsdeck/cli"
)

var (
	LogsCommand = cli.Command{
		Name:	"logs",
		ShortName: "l",
		Usage: "Print logs of current app",
		Action: func(c *cli.Context) {
			fmt.Println("logs")
		},
	}
)
