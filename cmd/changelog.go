package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/update"
)

var (
	changelogCommand = cli.Command{
		Name:     "changelog",
		Category: "CLI Internals",
		Usage:    "Show the scalingo CLI changelog from last version",
		Description: `Show the scalingo CLI changelog from last version
	Example
	  'scalingo changelog'`,
		Action: func(c *cli.Context) {
			err := update.ShowLastChangelog()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "changelog")
		},
	}
)
