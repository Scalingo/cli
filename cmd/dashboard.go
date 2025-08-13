package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
)

var (
	dashboardCommand = cli.Command{
		Name:     "dashboard",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Open app dashboard on default web browser",
		Description: CommandDescription{
			Description: "Open app dashboard on default web browser",
			Examples:    []string{"scalingo --app my-app dashboard"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(ctx, c, "dashboard")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			currentRegion := config.C.ScalingoRegion
			err := apps.Dashboard(currentApp, currentRegion)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "dashboard")
		},
	}
)
