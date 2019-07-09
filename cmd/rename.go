package cmd

import (
	"errors"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
)

var (
	renameCommand = cli.Command{
		Name:     "rename",
		Category: "App Management",
		Flags: []cli.Flag{
			appFlag,
			cli.StringFlag{
				Name:  "new-name",
				Value: "<new name>",
				Usage: "New name to give to the app",
			},
		},
		Usage:       "Rename an application",
		Description: "Rename an app\n  Example:\n    'scalingo rename --app my-app --new-name my-app-production'",
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			newName := c.String("new-name")
			if newName == "<new name>" {
				errorQuit(errors.New("--new-name flag should be defined"))
			}
			err := apps.Rename(currentApp, newName)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "rename")
		},
	}
)
