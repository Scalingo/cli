package cmd

import (
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/regions"
	"github.com/urfave/cli"
)

var (
	RegionsListCommand = cli.Command{
		Name:     "regions",
		Category: "Global",
		Usage:    "List available regions",
		Description: `
   Example
     'scalingo regions'`,
		Action: func(c *cli.Context) {
			err := regions.List()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "regions")
		},
	}
)
