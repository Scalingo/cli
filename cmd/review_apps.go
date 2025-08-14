package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/reviewapps"
)

var (
	reviewAppsShowCommand = cli.Command{
		Name:     "review-apps",
		Category: "Review Apps",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Show review apps of the parent application",
		Description: CommandDescription{
			Description: "Show review apps of the parent application",
			Examples:    []string{"scalingo --app my-app review-apps"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "review-apps")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := reviewapps.Show(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "review-apps")
		},
	}
)
