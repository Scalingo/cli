package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	logsCommand = cli.Command{
		Name:     "logs",
		Aliases:  []string{"l"},
		Category: "App Management",
		Usage:    "Get the logs of your applications",
		Description: CommandDescription{
			Description: "Get the logs of your applications",
			Examples: []string{
				"scalingo --app my-app logs -n 100                      # Get 100 lines",
				"scalingo --app my-app logs -f                          # Real-time logs",
				"scalingo --app my-app --addon addon_uuid logs          # Addon logs",
				"# Get lines with filter",
				"scalingo --app my-app logs -F web",
				"scalingo --app my-app logs -F web-1",
				"scalingo --app my-app logs --follow -F \"worker|clock\"",
			},
		}.Render(),
		Flags: []cli.Flag{&appFlag, &addonFlag,
			&cli.IntFlag{Name: "lines", Aliases: []string{"n"}, Value: 20, Usage: "Number of log lines to dump"},
			&cli.BoolFlag{Name: "follow", Aliases: []string{"f"}, Usage: "Stream logs of app, (as \"tail -f\")"},
			&cli.StringFlag{Name: "filter", Aliases: []string{"F"}, Usage: "Filter containers logs that will be displayed"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "logs")
				return nil
			}

			addonName := addonUUIDFromFlags(c, currentApp)

			var err error
			if addonName == "" {
				utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

				err = apps.Logs(c.Context, currentApp, c.Bool("f"), c.Int("n"), c.String("F"))
			} else {
				utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)

				err = db.Logs(c.Context, currentApp, addonName, db.LogsOpts{
					Follow: c.Bool("f"),
					Count:  c.Int("n"),
				})
			}

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "logs")
		},
	}
)
