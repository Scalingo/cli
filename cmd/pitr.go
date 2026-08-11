package cmd

import (
	"context"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/db/pitr"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v3"
)

var databasePITRRestore = cli.Command{
	Name:      "database-pitr-restore",
	Category:  "PITR management",
	Usage:     "Restore a database to a specific point in time",
	ArgsUsage: "restore-time",
	Flags: []cli.Flag{
		&appFlag,
		&addonFlag,
	},
	Description: CommandDescription{
		Description: "Restore a database to a specific point in time",
		Examples: []string{
			"scalingo --app my-app --addon my-addon database-pitr-restore 2026-07-18T23:00:00Z",
		},
	}.Render(),
	Action: func(ctx context.Context, c *cli.Command) error {
		currentResource, currentDatabase := detect.GetCurrentResourceAndDatabase(ctx, c)

		utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeDBs)
		addonName := currentDatabase
		if currentDatabase == "" {
			addonName = addonUUIDFromFlags(ctx, c, currentResource, true)
		}

		restoreTimeStr, err := time.Parse(time.RFC3339, c.Args().First())
		if err != nil {
			errorQuitWithHelpMessage(ctx, errors.Wrap(ctx, err, "invalid restore-time format"), c, "database-pitr-restore")
		}

		err = pitr.Restore(ctx, currentResource, addonName, restoreTimeStr)
		if err != nil {
			errorQuit(ctx, err)
		}

		return nil
	},
}
