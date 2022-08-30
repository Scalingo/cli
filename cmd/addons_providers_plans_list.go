package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/addon_providers"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	AddonProvidersPlansCommand = cli.Command{
		Name:        "addons-plans",
		Category:    "Addons - Global",
		Description: "List the plans for an addon.\n    Example:\n    scalingo addon-plans scalingo-mongodb",
		Usage:       "List plans",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "addons-plans")
				return nil
			}
			if err := addon_providers.Plans(c.Context, c.Args().First()); err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.AddonsPlansAutoComplete(c)
		},
	}
)
