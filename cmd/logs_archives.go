package cmd

import (
	"github.com/urfave/cli/v2"

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
		Flags: []cli.Flag{&appFlag, &addonFlag,
			&cli.IntFlag{Name: "page", Aliases: []string{"p"}, Usage: "Page number"},
		},
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "logs-archives")
				return nil
			}

			addonName := addonNameFromFlags(c)

			var err error
			if addonName == "" {
				utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

				err = apps.LogsArchives(c.Context, currentApp, c.Int("p"))
			} else {
				utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeDBs)

				err = db.LogsArchives(c.Context, currentApp, addonName, c.Int("p"))
			}

			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logs-archives")
		},
	}
)
