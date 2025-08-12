package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/update"
)

var (
	UpdateCommand = cli.Command{
		Name:        "update",
		Category:    "CLI Internals",
		Usage:       "Update 'scalingo' SDK client",
		Description: "Update 'scalingo' SDK client",
		Action: func(c *cli.Context) error {
			err := update.Check()
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "update")
		},
	}
)
