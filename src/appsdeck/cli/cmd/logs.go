package cmd

import (
	"appsdeck/cli/appdetect"
	"appsdeck/cli/apps"
	"github.com/Appsdeck/cli"
)

var (
	LogsCommand = cli.Command{
		Name:      "logs",
		ShortName: "l",
		Usage:     "Print logs of current app",
		Flags: []cli.Flag{
			cli.BoolFlag{"stream", "Stream logs of app, (as \"tail -f\")"},
		},
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if c.Bool("stream") {
				if err := apps.LogsStream(currentApp); err != nil {
					errorQuit(err)
				}
			} else if len(c.Args()) == 0 {
				if err := apps.Logs(currentApp); err != nil {
					errorQuit(err)
				}
			} else {
				cli.ShowCommandHelp(c, "logs")
			}
		},
	}
)
