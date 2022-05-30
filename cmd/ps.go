package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	psCommand = cli.Command{
		Name:     "ps",
		Category: "App Management",
		Usage:    "Display your application containers",
		Flags:    []cli.Flag{appFlag},
		Description: `Display your application containers
	Example
	  'scalingo --app my-app ps'`,
		Action: func(c *cli.Context) {
			currentApp := detect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "ps")
				return
			}

			err := apps.Ps(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "ps")
		},
	}
)
