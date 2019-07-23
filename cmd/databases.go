package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		Flags: []cli.Flag{appFlag, addonFlag, cli.StringFlag{
			Name:  "schedule-at",
			Usage: "Enable daily backups and schedule them at the specified hour of the day (in local time zone). It is also possible to specify the timezone to use.",
		}, cli.BoolFlag{
			Name:  "unschedule",
			Usage: "Disable the periodic backups",
		}},
		Description: `  Configure the periodic backups of a database:

Examples
 $ scalingo --app myapp --addon addon_uuid backups-config --schedule-at 3
 $ scalingo --app myapp --addon addon_uuid backups-config --schedule-at "3 Europe/Paris"
 $ scalingo --app myapp --addon addon_uuid backups-config --unschedule

		# See also 'addons' and 'backup-download'
		`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			addon := addonName(c)

			params := scalingo.PeriodicBackupsConfigParams{}
			scheduleAtFlag := c.String("schedule-at")
			disable := c.Bool("unschedule")
			if scheduleAtFlag != "" && disable {
				errorQuit(errors.New("You cannot use both --schedule-at and --unschedule at the same time"))
			}

			if disable {
				f := false
				params.Enabled = &f
			}
			if scheduleAtFlag != "" {
				t := true
				params.Enabled = &t
				scheduleAt, loc, err := parseScheduleAtFlag(scheduleAtFlag)
				if err != nil {
					errorQuit(err)
				}
				localTime := time.Date(1986, 7, 22, scheduleAt, 0, 0, 0, loc)
				hour := localTime.UTC().Hour()
				params.ScheduledAt = &hour
			}

			var database scalingo.Database
			var err error
			if disable || scheduleAtFlag != "" {
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

func parseScheduleAtFlag(flag string) (int, *time.Location, error) {
	scheduleAt, err := strconv.Atoi(flag)
	if err == nil {
		// In this case, the schedule-at flag equals a single number
		return scheduleAt, time.Local, nil
	}

	// From now on the schedule-at flag is a number and a timezone such as
	// "3 Europe/Paris"
	s := strings.Split(flag, " ")
	if len(s) < 2 {
		return -1, nil, errors.New("fail to parse the schedule-at flag")
	}
	scheduleAt, err = strconv.Atoi(s[0])
	if err != nil {
		return -1, nil, errors.New("fail to parse the schedule-at flag")
	}
	loc, err := time.LoadLocation(s[1])
	if err != nil {
		return -1, nil, fmt.Errorf("unknown timezone '%s'", s[1])
	}

	return scheduleAt, loc, nil
}
