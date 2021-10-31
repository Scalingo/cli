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
		Flags:    []cli.Flag{appFlag, 
			cli.StringFlag{Name: "region", Value: "", Usage: "Specify the region which contains the app"},},
		Usage:    "Open dashboard on default web browser",
		Description: `Open dashboard on default web browser:

	$ scalingo --region osc-fr1 --app my-app dashboard`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "dashboard")
				return
			}
			if c.String("region") == "" {
				cli.ShowCommandHelp(c, "dashboard")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := apps.Dashboard(currentApp, c.String("region"))
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "dashboard")
		},
	}
)
