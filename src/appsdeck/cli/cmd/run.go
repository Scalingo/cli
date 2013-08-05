package cmd

import (
	"fmt"
	"github.com/Appsdeck/cli"
)

var (
	RunCommand = cli.Command{
		Name:	"run",
		ShortName: "r",
		Usage: "Run command in current app context",
		Action: func(c *cli.Context) {
			fmt.Println("run")
		},
	}
)
