package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/regions"
)

var (
	RegionsListCommand = cli.Command{
		Name:     "regions",
		Category: "Global",
		Usage:    "List available regions",
		Description: `
   Example
     'scalingo regions'`,
		Action: func(c *cli.Context) error {
			err := regions.List()
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "regions")
		},
	}
)
