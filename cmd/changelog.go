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
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "since", Usage: "Show the changelog of every version released after this one"},
		},
		Description: CommandDescription{
			Description: "Show the Scalingo CLI changelog from last version, or of every version released after the one given with --since",
			Examples:    []string{"scalingo changelog", "scalingo changelog --since 1.43.1"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			err := update.ShowChangelogSince(ctx, c.String("since"))
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
