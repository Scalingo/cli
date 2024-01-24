package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/db"
	dbUsers "github.com/Scalingo/cli/db/users"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
)

var (
	databaseEnableFeature = cli.Command{
		Name:      "database-enable-feature",
		Category:  "Addons",
		Usage:     "Enable a togglable feature from a database",
		ArgsUsage: "feature-id",
		Flags: []cli.Flag{&appFlag, &addonFlag, &cli.BoolFlag{
			Name:  "synchronous",
			Usage: "Wait for the feature to be enabled synchronously",
		}},
		Description: CommandDescription{
			Description: "Enable a togglable feature from a database",
			Examples: []string{
				"scalingo --app myapp --addon addon-uuid database-enable-feature force-ssl",
				"scalingo --app myapp --addon addon-uuid database-enable-feature --synchronous force-ssl",
				"scalingo --app myapp --addon addon-uuid database-enable-feature publicly-available",
			},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
			addonName := addonUUIDFromFlags(c, currentApp, true)
			if c.NArg() != 1 {
				errorQuit(c.Context, errors.New("feature argument should be specified"))
			}
			feature := c.Args().First()
			err := db.EnableFeature(c, currentApp, addonName, feature)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}

	databaseDisableFeature = cli.Command{
		Name:      "database-disable-feature",
		Category:  "Addons",
		Usage:     "Enable a togglable feature from a database",
		ArgsUsage: "feature-id",
		Flags:     []cli.Flag{&appFlag, &addonFlag},
		Description: CommandDescription{
			Description: "Disable a togglable feature from a database",
			Examples: []string{
				"scalingo --app myapp --addon addon-uuid database-disable-feature force-ssl",
				"scalingo --app myapp --addon addon-uuid database-disable-feature publicly-available",
			},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
			addonName := addonUUIDFromFlags(c, currentApp, true)
			if c.NArg() != 1 {
				errorQuit(c.Context, errors.New("feature argument should be specified"))
			}
			feature := c.Args().First()
			err := db.DisableFeature(c.Context, currentApp, addonName, feature)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}

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
		Description: CommandDescription{
			Description: "Configure the periodic backups of a database",
			Examples: []string{
				"scalingo --app myapp --addon addon-uuid backups-config --schedule-at 3",
				"scalingo --app myapp --addon addon-uuid backups-config --schedule-at \"3 Europe/Paris\"",
				"scalingo --app myapp --addon addon-uuid backups-config --unschedule",
			},
			SeeAlso: []string{"addons", "backup-download"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
			addonName := addonUUIDFromFlags(c, currentApp, true)

			params := scalingo.DatabaseUpdatePeriodicBackupsConfigParams{}
			scheduleAtFlag := c.String("schedule-at")
			disable := c.Bool("unschedule")
			if scheduleAtFlag != "" && disable {
				errorQuit(c.Context, errors.New("you cannot use both --schedule-at and --unschedule at the same time"))
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
					errorQuit(c.Context, err)
				}
				localTime := time.Date(1986, 7, 22, scheduleAt, 0, 0, 0, loc)
				hour := localTime.UTC().Hour()
				params.ScheduledAt = &hour
			}

			if disable || scheduleAtFlag != "" {
				err := db.BackupsConfiguration(c.Context, currentApp, addonName, params)
				if err != nil {
					errorQuit(c.Context, err)
				}
			}
			return nil
		},
	}

	databaseListUsers = cli.Command{
		Name:     "database-users-list",
		Aliases:  []string{"database-list-users"},
		Category: "Addons",
		Usage:    "Print database's users",
		Flags:    []cli.Flag{&appFlag, &addonFlag},
		Description: CommandDescription{
			Description: `List the users of a database

Only available on ` + fmt.Sprintf("%s", dbUsers.SupportedAddons),
			Examples: []string{
				"scalingo --app myapp --addon addon-uuid database-users-list",
			},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
			addonName := addonUUIDFromFlags(c, currentApp, true)

			err := dbUsers.List(c.Context, currentApp, addonName)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}

	databaseDeleteUser = cli.Command{
		Name:      "database-users-delete",
		Aliases:   []string{"database-delete-user"},
		Category:  "Addons",
		ArgsUsage: "user",
		Usage:     "Delete a database's user",
		Flags:     []cli.Flag{&appFlag, &addonFlag},
		Description: CommandDescription{
			Description: `Delete the given user of a database

Only available on ` + fmt.Sprintf("%s", dbUsers.SupportedAddons),
			Examples: []string{
				"scalingo --app myapp --addon addon-uuid database-users-delete my_user",
			},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(c, "database-users-delete")
			}

			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
			addonName := addonUUIDFromFlags(c, currentApp, true)

			username := c.Args().First()

			err := dbUsers.DeleteUser(c.Context, currentApp, addonName, username)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}

	databaseCreateUser = cli.Command{
		Name:      "database-users-create",
		Aliases:   []string{"database-create-user"},
		Category:  "Addons",
		ArgsUsage: "user",
		Usage:     "Create new database user",
		Flags: []cli.Flag{
			&appFlag,
			&addonFlag,
			&cli.BoolFlag{Name: "read-only", Usage: "Create a user with read-only rights"},
		},
		Description: CommandDescription{
			Description: `Create new database user

Only available on ` + fmt.Sprintf("%s", dbUsers.SupportedAddons),
			Examples: []string{
				"scalingo --app myapp --addon addon-uuid database-users-create my_user",
				"scalingo --app myapp --addon addon-uuid database-users-create --read-only my_user",
			},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowCommandHelp(c, "database-users-create")
			}

			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)
			addonName := addonUUIDFromFlags(c, currentApp, true)

			username := c.Args().First()

			err := dbUsers.CreateUser(c.Context, currentApp, addonName, username, c.Bool("read-only"))
			if err != nil {
				errorQuit(c.Context, err)
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
