package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
)

var (
	LogsCommand = cli.Command{
		Name:      "logs",
		ShortName: "l",
		Category:  "App Management",
		Usage:     "Get the logs of your applications",
		Description: `Get the logs of your applications
   Example:
     Get 100 lines:          'scalingo --app my-app logs -n 100'
     Real-Time logs:         'scalingo --app my-app logs -f'
     Get lines with filter:  'scalingo --app my-app logs -n 100 -F web-1|worker-2'`,
		Flags: []cli.Flag{appFlag,
			cli.IntFlag{Name: "lines, n", Value: 20, Usage: "Number of log lines to dump", EnvVar: ""},
			cli.BoolFlag{Name: "follow, f", Usage: "Stream logs of app, (as \"tail -f\")", EnvVar: ""},
			cli.StringFlag{Name: "filter, F", Usage: "Filter process that will be displayed", EnvVar: ""},
		},
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) == 0 {
				if err := apps.Logs(currentApp, c.Bool("f"), c.Int("n"), c.String("F")); err != nil {
					errorQuit(err)
				}
			} else {
				cli.ShowCommandHelp(c, "logs")
			}
		},
	}
)
