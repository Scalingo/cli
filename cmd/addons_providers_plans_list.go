package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/addonproviders"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	AddonProvidersPlansCommand = cli.Command{
		Name:      "addons-plans",
		Category:  "Addons - Global",
		ArgsUsage: "addon-id",
		Description: CommandDescription{
			Description: "List the plans for an addon",
			Examples:    []string{"scalingo addon-plans scalingo-mongodb"},
		}.Render(),
		Usage: "List plans",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "addons-plans")
				return nil
			}
			if err := addonproviders.Plans(c.Context, c.Args().First()); err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.AddonsPlansAutoComplete(c)
		},
	}
)
