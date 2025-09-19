package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/dbng"
)

var (
	databasesList = cli.Command{
		Name:        "databases",
		Category:    "Databases",
		Usage:       "List the databases next generation that you own",
		Description: "List all the databases next generation of which you are an owner",
		Action: func(ctx context.Context, _ *cli.Command) error {
			err := dbng.List(ctx)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "databases")
		},
	}
)
