package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
)

var (
	dashboardCommand = cli.Command{
		Name:     "dashboard",
		Category: "App Management",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Open app dashboard on default web browser",
		Description: `Open app dashboard on default web browser:

	$ scalingo --app my-app dashboard`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "dashboard")
				return
			}

			currentApp := detect.CurrentApp(c)
			currentRegion := config.C.ScalingoRegion
			err := apps.Dashboard(currentApp, currentRegion)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "dashboard")
		},
	}
)
