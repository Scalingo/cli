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
		Usage:    "List available availability regions",
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

	RegionsSetCommand = cli.Command{
		Name:     "regions-set",
		Category: "Global",
		Usage:    "Configure the CLI to use a specific region",
		Description: `
   Example
     'scalingo regions-set agora-fr1'

	 Can also be configured using the environment variable
	   SC_REGION=agora-fr1`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "regions-set")
			} else {
				err := regions.Set(c.Args()[0])
				if err != nil {
					errorQuit(err)
				}
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "regions-set")
		},
	}
)
