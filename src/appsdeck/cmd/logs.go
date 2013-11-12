package cmd

import (
	"appsdeck/appdetect"
	"appsdeck/apps"
	"github.com/codegangsta/cli"
)

var (
	LogsCommand = cli.Command{
		Name:      "logs",
		ShortName: "l",
		Description: "Get the logs of an application",
		Usage:     "[-n <nblines> | --stream]",
		Flags: []cli.Flag{
			cli.IntFlag{"n", 20, "Number of log lines to dump"},
			cli.BoolFlag{"stream", "Stream logs of app, (as \"tail -f\")"},
		},
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if c.Bool("stream") {
				if err := apps.LogsStream(currentApp); err != nil {
					errorQuit(err)
				}
			} else if len(c.Args()) == 0 || len(c.Args()) == 1 && c.Int("n") != 0 {
				if err := apps.Logs(currentApp, c.Int("n")); err != nil {
					errorQuit(err)
				}
			} else {
				cli.ShowCommandHelp(c, "logs")
			}
		},
	}
)
