package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v2"
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
			currentResource := detect.GetCurrentResource(ctx, c)
			isDB, err := utils.IsResourceDatabase(ctx, currentResource)
			if err != nil && !errors.Is(err, utils.ErrResourceNotFound) {
				errorQuit(ctx, err)
			}
			if isDB {
				io.Error("It is currently impossible to rename a database.")
				return nil
			}

			newName := c.String("new-name")

			utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeContainers)

			err = apps.Rename(ctx, currentResource, newName)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "rename")
		},
	}
)
