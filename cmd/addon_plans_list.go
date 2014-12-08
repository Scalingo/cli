package cmd

import (
	"github.com/Scalingo/cli/addons"
	"github.com/codegangsta/cli"
)

var (
	AddonPlansCommand = cli.Command{
		Name:        "addon-plans",
		Description: "List the plans for an addon.\n    Example:\n    scalingo addon-plans scalingo-mongo",
		Usage:       "List plans",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "addon-plans")
				return
			}
			if err := addons.Plans(c.Args()[0]); err != nil {
				errorQuit(err)
			}
		},
	}
)
