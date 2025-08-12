package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/regions"
)

var (
	RegionsListCommand = cli.Command{
		Name:        "regions",
		Category:    "Global",
		Usage:       "List available regions",
		Description: "List available regions",
		Action: func(c *cli.Context) error {
			err := regions.List(c.Context)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "regions")
		},
	}
)
