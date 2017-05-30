package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	LogsArchivesCommand = cli.Command{
		Name:      "logs-archives",
		ShortName: "la",
		Category:  "App Management",
		Usage:     "Get the logs archives of your applications",
		Description: `Get the logs archives of your applications
   Example:
     Get most recents archives : 'scalingo --app my-app logs-archives'
     Get a specific page : 'scalingo --app my-app logs-archives -c CURSOR'`,
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "cursor, c", Usage: "Cursor to the next page", EnvVar: ""},
		},
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) == 0 {
				if err := apps.LogsArchives(currentApp, c.String("c")); err != nil {
					errorQuit(err)
				}
			} else {
				cli.ShowCommandHelp(c, "logs")
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logs-archives")
		},
	}
)
