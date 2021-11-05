package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
)

var (
	logsArchivesCommand = cli.Command{
		Name:      "logs-archives",
		ShortName: "la",
		Category:  "App Management",
		Usage:     "Get the logs archives of your applications and databases",
		Description: `Get the logs archives of your applications and databases
   Examples:
     Get most recents archives: 'scalingo --app my-app logs-archives'
     Get a specific page:       'scalingo --app my-app logs-archives -p 5'
	   Addon logs archives:       'scalingo --app my-app logs-archives --addon addon-id'`,
		Flags: []cli.Flag{appFlag, addonFlag,
			cli.IntFlag{Name: "page, p", Usage: "Page number", EnvVar: ""},
		},
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "logs-archives")
				return
			}

			addonName := addonNameFromFlags(c)

			var err error
			if addonName == "" {
				err = apps.LogsArchives(currentApp, c.Int("p"))
			} else {
				err = db.LogsArchives(currentApp, addonName, c.Int("p"))
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logs-archives")
		},
	}
)
