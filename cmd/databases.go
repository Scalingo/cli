package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v5"
)

var (
	databaseBackupsConfig = cli.Command{
		Name:     "backups-config",
		Category: "Addons",
		Usage:    "Configure the periodic backups of a database",
		Flags: []cli.Flag{&appFlag, &addonFlag, &cli.StringFlag{
			Name:  "schedule-at",
			Usage: "Enable daily backups and schedule them at the specified hour of the day (in local time zone). It is also possible to specify the timezone to use.",
		}, &cli.BoolFlag{
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
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			addonName := addonNameFromFlags(c)
			if addonName == "" {
				fmt.Println("Unable to find the addon name, please use --addon flag.")
				os.Exit(1)
			}

			params := scalingo.PeriodicBackupsConfigParams{}
			scheduleAtFlag := c.String("schedule-at")
			disable := c.Bool("unschedule")
			if scheduleAtFlag != "" && disable {
				errorQuit(errors.New("You cannot use both --schedule-at and --unschedule at the same time"))
			}
			database, err := db.Show(c.Context, currentApp, addonName)
			if err != nil {
				errorQuit(err)
			}
			if scheduleAtFlag != "" && len(database.PeriodicBackupsScheduledAt) > 1 {
				msg := "Your database is backed up multiple times a day at " +
					formatScheduledAt(database.PeriodicBackupsScheduledAt) +
					". Please ask the support to update the frequency of these backups."
				errorQuit(errors.New(msg))
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

			if disable || scheduleAtFlag != "" {
				database, err = db.BackupsConfiguration(c.Context, currentApp, addonName, params)
				if err != nil {
					errorQuit(err)
				}
			}
			if database.PeriodicBackupsEnabled {
				io.Statusf("Periodic backups will be done daily at %s\n", formatScheduledAt(database.PeriodicBackupsScheduledAt))
			} else {
				io.Status("Periodic backups are disabled")
			}
			return nil
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

func formatScheduledAt(hours []int) string {
	hoursStr := make([]string, len(hours))
	for i, h := range hours {
		hUTC := time.Date(1986, 7, 22, h, 0, 0, 0, time.UTC)
		hLocal := hUTC.In(time.Local)
		hoursStr[i] = strconv.Itoa(hLocal.Hour())
	}

	tz, _ := time.Now().In(time.Local).Zone()
	return fmt.Sprintf("%s:00 %s", strings.Join(hoursStr, ":00, "), tz)
}
