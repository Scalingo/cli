package cmd

import (
	"context"

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
		Flags: []cli.Flag{&appFlag, &addonFlag, databaseFlag(),
			&cli.IntFlag{Name: "lines", Aliases: []string{"n"}, Value: 20, Usage: "Number of log lines to dump"},
			&cli.BoolFlag{Name: "follow", Aliases: []string{"f"}, Usage: "Stream logs of app, (as \"tail -f\")"},
			&cli.StringFlag{Name: "filter", Aliases: []string{"F"}, Usage: "Filter containers logs that will be displayed"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource := detect.GetCurrentResource(ctx, c)
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "logs")
				return nil
			}

			addonID := addonUUIDFromFlags(ctx, c, currentResource)

			var err error
			if addonID == "" {
				utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeContainers)

				err = apps.Logs(ctx, currentResource, c.Bool("f"), c.Int("n"), c.String("F"))
			} else {
				utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeDBs)

				err = db.Logs(ctx, currentResource, addonID, db.LogsOpts{
					Follow: c.Bool("f"),
					Count:  c.Int("n"),
				})
			}

			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "logs")
		},
	}
)
