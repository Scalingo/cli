package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
)

var (
	logsCommand = cli.Command{
		Name:     "logs",
		Aliases:  []string{"l"},
		Category: "App Management",
		Usage:    "Get the logs of your applications",
		Description: `Get the logs of your applications
   Example:
     Get 100 lines:          'scalingo --app my-app logs -n 100'
     Real-Time logs:         'scalingo --app my-app logs -f'
     Addon logs:             'scalingo --app my-app --addon addon_uuid logs'
     Get lines with filter:
       'scalingo --app my-app logs -F web'
       'scalingo --app my-app logs -F web-1'
       'scalingo --app my-app logs -F router'
       'scalingo --app my-app logs -F one-off'
       'scalingo --app my-app logs -F one-off-1'
       'scalingo --app my-app logs --follow -F "web|worker"'`,
		Flags: []cli.Flag{&appFlag, &addonFlag,
			&cli.IntFlag{Name: "lines", Aliases: []string{"n"}, Value: 20, Usage: "Number of log lines to dump"},
			&cli.BoolFlag{Name: "follow", Aliases: []string{"f"}, Usage: "Stream logs of app, (as \"tail -f\")"},
			&cli.StringFlag{Name: "filter", Aliases: []string{"F"}, Usage: "Filter containers logs that will be displayed"},
		},
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "logs")
				return nil
			}

			addonName := addonNameFromFlags(c)

			var err error
			if addonName == "" {
				err = apps.Logs(c.Context, currentApp, c.Bool("f"), c.Int("n"), c.String("F"))
			} else {
				err = db.Logs(c.Context, currentApp, addonName, db.LogsOpts{
					Follow: c.Bool("f"),
					Count:  c.Int("n"),
				})
			}

			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logs")
		},
	}
)
