package cmd

import (
	"github.com/urfave/cli"
	"github.com/Scalingo/cli/addon_providers"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	AddonProvidersPlansCommand = cli.Command{
		Name:        "addons-plans",
		Category:    "Addons - Global",
		Description: "List the plans for an addon.\n    Example:\n    scalingo addon-plans scalingo-mongo",
		Usage:       "List plans",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "addon-plans")
				return
			}
			if err := addon_providers.Plans(c.Args()[0]); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.AddonsPlansAutoComplete(c)
		},
	}
)
