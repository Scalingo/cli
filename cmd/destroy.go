package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	DestroyCommand = cli.Command{
		Name:     "destroy",
		Category: "App Management",
		Flags: []cli.Flag{&appFlag,
			&cli.BoolFlag{Name: "force", Usage: "Force destroy without asking for a confirmation /!\\"},
		},
		Usage: "Destroy an app /!\\",
		Description: `Destroy an app /!\\ It is not reversible	

Example:
 $ scalingo destroy my-app
 $ scalingo -a my-app destroy --force'
`,
		Action: func(c *cli.Context) error {
			var currentApp string

			if c.Args().Len() > 1 {
				cli.ShowCommandHelp(c, "destroy")
			} else {
				if c.Args().Len() != 0 {
					currentApp = c.Args().First()
				} else {
					currentApp = detect.CurrentApp(c)
				}

				err := apps.Destroy(c.Context, currentApp, c.Bool("force"))
				if err != nil {
					errorQuit(err)
				}
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "destroy")
		},
	}
)
