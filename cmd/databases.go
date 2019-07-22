package cmd

import (
	"errors"
	"time"

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
			Name:  "schedule-at",
			Usage: "Enable daily backups and schedule them at the specified hour of the day (in local time zone)",
		}, cli.BoolFlag{
			Name:  "disable-scheduling",
			Usage: "Disable the periodic backups",
		}},
		Description: `  Configure the periodic backups of a database:
		$ scalingo --app myapp --addon addon_uuid backups-config --schedule-at 3

		# See also 'addons' and 'backup-download'
		`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)

			params := scalingo.PeriodicBackupsConfigParams{}
			scheduleAt := c.Int("schedule-at")
			disable := c.Bool("disable-scheduling")
			if scheduleAt != 0 && disable {
				errorQuit(errors.New("You cannot use both --schedule-at and --disable-scheduling at the same time"))
			}

			if disable {
				f := false
				params.Enabled = &f
			}
			if scheduleAt != 0 {
				t := true
				params.Enabled = &t
				localTime := time.Date(1986, 7, 22, scheduleAt, 0, 0, 0, time.Local)
				hour := localTime.UTC().Hour()
				params.ScheduledAt = &hour
			}

			var database scalingo.Database
			var err error
			if disable || scheduleAt != 0 {
				database, err = db.BackupsConfiguration(currentApp, addon, params)
				if err != nil {
					errorQuit(err)
				}
			} else {
				database, err = db.Show(currentApp, addon)
			}
			if database.PeriodicBackupsEnabled {
				scheduledAtUTC := time.Date(1986, 7, 22, database.PeriodicBackupsScheduledAt, 0, 0, 0, time.UTC)
				scheduledAt := scheduledAtUTC.In(time.Local)
				zone, _ := time.Now().In(time.Local).Zone()
				io.Statusf("Periodic backups will be done daily at %d:00 %s\n", scheduledAt.Hour(), zone)
			} else {
				io.Status("Periodic backups are disabled")
			}
		},
	}
)
