package cmd

import (
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/update"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	UpdateCommand = cli.Command{
		Name:     "update",
		Category: "CLI Internals",
		Usage:    "Update 'scalingo' client",
		Description: `Update 'scalingo' client
   Example
     'scalingo update'`,
		Action: func(c *cli.Context) {
			err := update.Check()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "update")
		},
	}
)
