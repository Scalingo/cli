package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
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
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			newName := c.String("new-name")

			err := apps.Rename(c.Context, currentApp, newName)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "rename")
		},
	}
)
