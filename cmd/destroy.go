package cmd

import (
	"github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	DestroyCommand = cli.Command{
		Name:        "destroy",
		Category:    "Global",
		Flags:       []cli.Flag{appFlag},
		Usage:       "Destroy an app /!\\",
		Description: "Destroy an app /!\\ It is not reversible\n  Example:\n    'scalingo destroy my-app'\n    'scalingo -a my-app destroy'",
		Action: func(c *cli.Context) {
			var currentApp string

			if len(c.Args()) > 1 {
				cli.ShowCommandHelp(c, "destroy")
			} else {
				if len(c.Args()) != 0 {
					currentApp = c.Args()[0]
				} else {
					currentApp = appdetect.CurrentApp(c)
				}

				err := apps.Destroy(currentApp)
				if err != nil {
					errorQuit(err)
				}
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "destroy")
		},
	}
)
