package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	DestroyCommand = cli.Command{
		Name:     "destroy",
		Category: "App Management",
		Flags: []cli.Flag{appFlag,
			cli.BoolFlag{Name: "force", Usage: "Force destroy without asking for a confirmation /!\\", EnvVar: ""},
		},
		Usage: "Destroy an app /!\\",
		Description: "Destroy an app /!\\ It is not reversible\n	Example:\n    'scalingo destroy my-app'\n    'scalingo -a my-app destroy --force'\n	",
		Action: func(c *cli.Context) {
			var currentApp string

			if len(c.Args()) > 1 {
				cli.ShowCommandHelp(c, "destroy")
			} else {
				if len(c.Args()) != 0 {
					currentApp = c.Args()[0]
				} else {
					currentApp = detect.CurrentApp(c)
				}

				err := apps.Destroy(currentApp, c.Bool("force"))
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
