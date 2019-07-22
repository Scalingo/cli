package cmd

import (
	"errors"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/urfave/cli"
)

var (
	databaseBackupsConfig = cli.Command{
		Name:     "backups-config",
		Category: "Addons",
		Usage:    "Configure the periodic backups of a database",
		Flags: []cli.Flag{appFlag, addonFlag, cli.IntFlag{
			Name:  "scheduled-at",
			Usage: "Hour of the day of the periodic backups (UTC)",
		}, cli.BoolFlag{
			Name:  "disable",
			Usage: "Disable the periodic backups",
		}, cli.BoolFlag{
			Name:  "enable",
			Usage: "Enable the periodic backups",
		}},
		Description: `  Configure the periodic backups of a databas:
		$ scalingo --app myapp --addon addon_uuid backups-config

		# See also 'addons' and 'backup-download'
		`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)
			params := scalingo.PeriodicBackupsConfigParams{}

			enable := c.Bool("enable")
			disable := c.Bool("disable")
			if enable && disable {
				errorQuit(errors.New("You cannot use both --enable and --disable at the same time"))
			}
			if enable {
				t := true
				params.Enabled = &t
			}
			if disable {
				f := false
				params.Enabled = &f
			}

			scheduledAt := c.Int("scheduled-at")
			if scheduledAt != 0 {
				params.ScheduledAt = &scheduledAt
			}

			db, err := db.BackupsConfiguration(currentApp, addon, params)
			if err != nil {
				errorQuit(err)
			}
			if db.PeriodicBackupsEnabled {
				io.Statusf("Periodic backups will be done daily at %d:00 UTC\n", db.PeriodicBackupsScheduledAt)
			} else {
				io.Statusf("Periodic backups are disabled")
			}
		},
	}
)
