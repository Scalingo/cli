package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/db/maintenance"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo/v6"
)

var databaseMaintenanceList = cli.Command{
	Name:     "database-maintenance-list",
	Category: "Addons",
	Usage:    "List database maintenance",
	Flags: []cli.Flag{
		&appFlag,
		&addonFlag,
		&cli.IntFlag{Name: "page", Usage: "Page to display", Value: 1},
		&cli.IntFlag{Name: "per-page", Usage: "Number of deployments to display", Value: 20},
	},
	Description: CommandDescription{
		Description: "List database maintenance",
		Examples:    []string{},
	}.Render(),

	Action: func(c *cli.Context) error {
		currentApp := detect.CurrentApp(c)
		utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
		addonName := addonNameFromFlags(c, true)

		err := maintenance.List(c.Context, currentApp, addonName, scalingo.PaginationOpts{
			Page:    c.Int("page"),
			PerPage: c.Int("per-page"),
		})
		if err != nil {
			errorQuit(err)
		}
		return nil
	},
}
