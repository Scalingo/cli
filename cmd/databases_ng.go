package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/dbng"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	databasesList = cli.Command{
		Name:     "databases",
		Category: "Databases",
		Usage:    "List the databases next generation that you own",
		Description: CommandDescription{
			Description: "List all the databases next generation of which you are an owner",
			SeeAlso:     []string{"database-info"},
		}.Render(),
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

	databaseShow = cli.Command{
		Name:     "database-info",
		Category: "Databases",
		Usage:    "View database next generation",
		Description: CommandDescription{
			Description: "View database next generation detailed informations",
			Examples:    []string{"scalingo database-info database_id"},
			SeeAlso:     []string{"databases"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(ctx, c, "database-info")
			}

			utils.CheckForConsent(ctx, currentApp)

			err := dbng.Show(ctx, c.Args().First())
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-info")
		},
	}
)
