package cmd

import (
	"context"

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
		Action: func(ctx context.Context, c *cli.Command) error {
			err := update.Check()
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "update")
		},
	}
)
