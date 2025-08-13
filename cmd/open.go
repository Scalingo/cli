package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
)

var (
	openCommand = cli.Command{
		Name:     "open",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Open app on default web browser",
		Description: CommandDescription{
			Description: "Open app on default web browser",
			Examples:    []string{"scalingo --app my-app open"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "open")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			currentRegion := config.C.ScalingoRegion

			err := apps.Open(currentApp, currentRegion)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "open")
		},
	}
)
