package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/update"
)

var (
	changelogCommand = cli.Command{
		Name:     "changelog",
		Category: "CLI Internals",
		Usage:    "Show the Scalingo CLI changelog from last version",
		Description: CommandDescription{
			Description: "Show the Scalingo CLI changelog from last version",
			Examples:    []string{"scalingo changelog"},
		}.Render(),

		Action: func(c *cli.Context) error {
			err := update.ShowLastChangelog()
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "changelog")
		},
	}
)
