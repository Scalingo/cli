package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	dashboardCommand = cli.Command{
		Name:     "dashboard",
		Category: "App Management",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Open app dashboard on default web browser",
		Description: `Open app dashboard on default web browser:

	$ scalingo --region osc-fr1 --app my-app dashboard`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "dashboard")
				return
			}

			if c.GlobalString("region") == "" {
				cli.ShowCommandHelp(c, "dashboard")
			}

			currentApp := appdetect.CurrentApp(c)
			err := apps.Dashboard(currentApp, c.GlobalString("region"))
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "dashboard")
		},
	}
)
