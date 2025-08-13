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
	renameCommand = cli.Command{
		Name:     "rename",
		Category: "App Management",
		Flags: []cli.Flag{
			&appFlag,
			&cli.StringFlag{Name: "new-name", Usage: "New name to give to the app", Required: true},
		},
		Usage: "Rename an application",
		Description: CommandDescription{
			Description: "Rename an application",
			Examples:    []string{"scalingo rename --app my-app --new-name my-app-production"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			newName := c.String("new-name")

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			err := apps.Rename(ctx, currentApp, newName)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "rename")
		},
	}
)
