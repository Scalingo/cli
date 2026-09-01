package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	valkeyConsoleCommand = cli.Command{
		Name:     "valkey-console",
		Category: categoryDatabases,
		Usage:    "Run an interactive console with your Valkey addon",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{
				Name: flagSizeName, Aliases: []string{"s"}, Value: "",
				Usage: flagSizeUsage,
			},
			&cli.StringFlag{
				Name: flagEnvName, Aliases: []string{"e"}, Value: "",
				Usage: flagEnvUsage,
			},
		},
		Description: CommandDescription{
			Description: `Run an interactive console with your Valkey addon.

The --size flag makes it easy to specify the size of the container executing
the Valkey console. Each container size has different price and performance.
You can read more about container sizes here:
https://doc.scalingo.com/platform/internals/container-sizes`,
			Examples: []string{
				"scalingo --app my-app valkey-console",
				"scalingo --app my-app valkey-console --size L",
				"scalingo --app my-app valkey-console --env MY_VALKEY_URL",
			},
			SeeAlso: []string{mySQLConsole, pgSQLConsole},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "valkey-console")
				return nil
			}
			currentApp := detect.CurrentApp(ctx, c)

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeDBs)

			err := db.ValkeyConsole(ctx, db.ValkeyConsoleOpts{
				App:          currentApp,
				Size:         c.String("s"),
				VariableName: c.String("e"),
			})
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "valkey-console")
		},
	}
)
