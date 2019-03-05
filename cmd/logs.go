package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/urfave/cli"
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
     Get lines with filter:
       'scalingo --app my-app logs -F web'
       'scalingo --app my-app logs -F web-1'
       'scalingo --app my-app logs --follow -F "worker|clock"'`,
		Flags: []cli.Flag{appFlag,
			cli.IntFlag{Name: "lines, n", Value: 20, Usage: "Number of log lines to dump", EnvVar: ""},
			cli.BoolFlag{Name: "follow, f", Usage: "Stream logs of app, (as \"tail -f\")", EnvVar: ""},
			cli.StringFlag{Name: "filter, F", Usage: "Filter containers logs that will be displayed", EnvVar: ""},
		},
		Before: AuthenticateHook,
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
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logs")
		},
	}

	AddonsLogsCommand = cli.Command{
		Name:     "addon-logs",
		Category: "Addons",
		Usage:    "Get the logs of your addons",
		Description: `Get the logs of your addons
   Example:
     Get 100 lines:          'scalingo --app my-app --addon addon_uuid logs -n 100'
     Real-Time logs:         'scalingo --app my-app --addon addon_uuid logs -f'
		 `,
		Flags: []cli.Flag{appFlag,
			cli.IntFlag{Name: "lines, n", Value: 20, Usage: "Number of log lines to dump", EnvVar: ""},
			cli.BoolFlag{Name: "follow, f", Usage: "Stream logs of addon, (as \"tail -f\")", EnvVar: ""},
		},
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "logs")
				return
			}
			currentApp := appdetect.CurrentApp(c)
			currentAddon := addonName(c)
			opts := db.LogsOpts{
				Follow: c.Bool("follow"),
				Count:  c.Int("lines"),
			}
			err := db.Logs(currentApp, currentAddon, opts)
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
