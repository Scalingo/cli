package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/addon_providers"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	AddonProvidersListCommand = cli.Command{
		Name:        "addons-list",
		Category:    "Addons - Global",
		Description: "List all addons you can add to your app",
		Usage:       "List all addons",
		Action: func(c *cli.Context) error {
			if err := addon_providers.List(c.Context); err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-list")
		},
	}
)
