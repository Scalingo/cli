package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	DestroyCommand = cli.Command{
		Name:     "destroy",
		Category: "App Management",
		Flags: []cli.Flag{&appFlag,
			&cli.BoolFlag{Name: "force", Usage: "Force destroy without asking for a confirmation /!\\"},
		},
		Usage:     "Destroy an app /!\\",
		ArgsUsage: "<--app app-id | app-id>",
		Description: CommandDescription{
			Description: "Destroy an app /!\\ It is not reversible",
			Examples: []string{
				"scalingo destroy my-app",
				"scalingo --app my-app destroy --force",
			},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			var currentApp string

			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "destroy")
			} else {
				if c.Args().Len() != 0 {
					currentApp = c.Args().First()
				} else {
					currentApp = detect.CurrentApp(c)
				}

				utils.CheckForConsent(ctx, currentApp)
				err := apps.Destroy(ctx, currentApp, c.Bool("force"))
				if err != nil {
					errorQuit(ctx, err)
				}
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "destroy")
		},
	}
)
