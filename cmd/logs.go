package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/urfave/cli"
)

var (
	logsCommand = cli.Command{
		Name:      "logs",
		ShortName: "l",
		Category:  "App Management",
		Usage:     "Get the logs of your applications",
		Description: `Get the logs of your applications
   Example:
     Get 100 lines:          'scalingo --app my-app logs -n 100'
     Real-Time logs:         'scalingo --app my-app logs -f'
     Addon logs:             'scalingo --app my-app --addon addon_uuid logs'
     Get lines with filter:
       'scalingo --app my-app logs -F web'
       'scalingo --app my-app logs -F web-1'
       'scalingo --app my-app logs --follow -F "worker|clock"'`,
		Flags: []cli.Flag{appFlag, addonFlag,
			cli.IntFlag{Name: "lines, n", Value: 20, Usage: "Number of log lines to dump", EnvVar: ""},
			cli.BoolFlag{Name: "follow, f", Usage: "Stream logs of app, (as \"tail -f\")", EnvVar: ""},
			cli.StringFlag{Name: "filter, F", Usage: "Filter containers logs that will be displayed", EnvVar: ""},
		},
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "logs")
				return
			}
			var addonName string
			if c.GlobalString("addon") != "<addon_id>" {
				addonName = c.GlobalString("addon")
			} else if c.String("addon") != "<addon_id>" {
				addonName = c.String("addon")
			}
			var err error
			if addonName == "" {
				err = apps.Logs(currentApp, c.Bool("f"), c.Int("n"), c.String("F"))
			} else {
				err = db.Logs(currentApp, addonName, db.LogsOpts{
					Follow: c.Bool("f"),
					Count:  c.Int("n"),
				})

			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logs")
		},
	}
)
