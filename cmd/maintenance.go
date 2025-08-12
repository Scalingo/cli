package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/db/maintenance"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

var databaseMaintenanceList = cli.Command{
	Name:     "database-maintenance-list",
	Category: "Addons Maintenance",
	Usage:    "List the past and future maintenance on the given database",
	Flags: []cli.Flag{
		&appFlag,
		&addonFlag,
		&cli.IntFlag{Name: "page", Usage: "Page to display", Value: 1},
		&cli.IntFlag{Name: "per-page", Usage: "Number of deployments to display", Value: 20},
	},
	Description: CommandDescription{
		Description: "List database maintenance",
		Examples: []string{
			"scalingo --app my-app --addon addon-uuid database-maintenance-list",
			"scalingo --app my-app --addon addon-uuid database-maintenance-list --per-page 20 --page 5",
		},
	}.Render(),

	Action: func(c *cli.Context) error {
		currentApp := detect.CurrentApp(c)
		utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
		addonName := addonUUIDFromFlags(c, currentApp, true)

		err := maintenance.List(c.Context, currentApp, addonName, scalingo.PaginationOpts{
			Page:    c.Int("page"),
			PerPage: c.Int("per-page"),
		})
		if err != nil {
			errorQuit(c.Context, err)
		}
		return nil
	},
}

var databaseMaintenanceInfo = cli.Command{
	Name:     "database-maintenance-info",
	Category: "Addons",
	Usage:    "Show a database maintenance",
	Flags: []cli.Flag{
		&appFlag,
		&addonFlag,
		&cli.StringFlag{Name: "maintenance", Usage: "ID of the maintenance"},
	},
	Description: CommandDescription{
		Description: "Show a database maintenance",
		Examples: []string{
			"scalingo --app my-app --addon ad-9be0fc04-bee6-4981-a403-a9ddbee7bd1f database-maintenance-info 64a56b51a8acb50065b73ec8",
		},
	}.Render(),
	Action: func(c *cli.Context) error {
		currentApp := detect.CurrentApp(c)
		if c.Args().Len() != 1 {
			err := cli.ShowCommandHelp(c, "database-maintenance-info")
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		}

		utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
		addonName := addonUUIDFromFlags(c, currentApp, true)
		maintenanceID := c.Args().First()

		err := maintenance.Info(c.Context, currentApp, addonName, maintenanceID)
		if err != nil {
			errorQuit(c.Context, err)
		}
		return nil
	},
}
