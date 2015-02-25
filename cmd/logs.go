package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	LogsCommand = cli.Command{
		Name:      "logs",
		ShortName: "l",
		Category:  "App Management",
		Usage:     "Get the logs of your applications",
		Description: `Get the logs of your applications
   Example:
     Get 100 lines:  'scalingo --app my-app logs -n 100'
     Real-Time logs: 'scalingo --app my-app logs -f'`,
		Flags: []cli.Flag{appFlag,
			cli.IntFlag{Name: "lines, n", Value: 20, Usage: "Number of log lines to dump", EnvVar: ""},
			cli.BoolFlag{Name: "follow, f", Usage: "Stream logs of app, (as \"tail -f\")", EnvVar: ""},
		},
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
