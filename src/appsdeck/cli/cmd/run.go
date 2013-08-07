package cmd

import (
	"fmt"
	"github.com/Appsdeck/cli"
)

var (
	RunCommand = cli.Command{
		Name:      "run",
		ShortName: "r",
		Usage:     "Run command in current app context",
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if err := apps.Run(currentApp, c.Args); err != nil {
				errorQuit(err)
			} else {
				cli.ShowCommandHelp(c, "logs")
			}
		},
	}
)
