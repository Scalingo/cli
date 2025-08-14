package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	psCommand = cli.Command{
		Name:     "ps",
		Category: "App Management",
		Usage:    "Display your application containers",
		Flags:    []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Display your application containers",
			Examples:    []string{"scalingo --app my-app ps"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "ps")
				return nil
			}

			err := apps.Ps(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "ps")
		},
	}
)
