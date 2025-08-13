package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

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

		Action: func(ctx context.Context, _ *cli.Command) error {
			err := update.ShowLastChangelog()
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "changelog")
		},
	}
)
