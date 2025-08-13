package cmd

import (
	"context"

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
		Action: func(ctx context.Context, _ *cli.Command) error {
			err := regions.List(ctx)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "regions")
		},
	}
)
