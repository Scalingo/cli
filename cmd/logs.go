package cmd

import (
	"github.com/Appsdeck/appsdeck/appdetect"
	"github.com/Appsdeck/appsdeck/apps"
	"github.com/Appsdeck/appsdeck/auth"
	"github.com/codegangsta/cli"
)

var (
	LogsCommand = cli.Command{
		Name:      "logs",
		ShortName: "l",
		Usage:     "Get the logs of your applications",
		Description: `Get the logs of your applications
   Example:
     Get 100 lines:  'appsdeck --app my-app logs -n 100'
     Real-Time logs: 'appsdeck --app my-app logs -f'`,
		Flags: []cli.Flag{
			cli.IntFlag{"lines, n", 20, "Number of log lines to dump", ""},
			cli.BoolFlag{"follow, f", "Stream logs of app, (as \"tail -f\")", ""},
		},
		Action: func(c *cli.Context) {
			auth.InitAuth()
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if len(c.Args()) == 0 || len(c.Args()) == 2 && c.Int("n") != 0 {
				if err := apps.Logs(currentApp, c.Bool("f"), c.Int("n")); err != nil {
					errorQuit(err)
				}
			} else {
				cli.ShowCommandHelp(c, "logs")
			}
		},
	}
)
