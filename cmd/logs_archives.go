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
	logsArchivesCommand = cli.Command{
		Name:     "logs-archives",
		Aliases:  []string{"la"},
		Category: "App Management",
		Usage:    "Get the logs archives of your applications and databases",
		Description: CommandDescription{
			Description: "Get the logs archives of your applications and databases",
			Examples: []string{
				"scalingo --app my-app logs-archives                   # Get the most recent archives",
				"scalingo --app my-app logs-archives -p 5              # Get a specific page",
				"scalingo --app my-app logs-archives --addon addon-id  # Addon logs archives",
			},
		}.Render(),
		Flags: []cli.Flag{&appFlag, &addonFlag, databaseFlag(),
			&cli.IntFlag{Name: "page", Aliases: []string{"p"}, Usage: "Page number"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource, currentDatabase := detect.GetCurrentResourceAndDatabase(ctx, c)
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "logs-archives")
				return nil
			}

			addonID := currentDatabase
			if currentDatabase == "" {
				addonID = addonUUIDFromFlags(ctx, c, currentResource)
			}

			var err error
			if addonID == "" {
				utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeContainers)

				err = apps.LogsArchives(ctx, currentResource, c.Int("p"))
			} else {
				utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeDBs)

				err = db.LogsArchives(ctx, currentResource, addonID, c.Int("p"))
			}

			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "logs-archives")
		},
	}
)
