package cmd

import (
	"github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	RestartCommand = cli.Command{
		Name:     "restart",
		Category: "App Management",
		Usage:    "Restart processes of your app",
		Flags:    []cli.Flag{appFlag, cli.BoolFlag{Name: "synchronous, s", Usage: "Do the restart synchronously", EnvVar: ""}},
		Description: `Restart one or several process or your application:
	Example
	  ## Restart all the processes
	  scalingo --app my-app restart
		## Restart all the web processes
	  scalingo --app my-app restart web
		## Restart a specific container
	  scalingo --app my-app restart web-1`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if err := apps.Restart(currentApp, c.Bool("s"), c.Args()); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "restart")
			autocomplete.RestartAutoComplete(c)
		},
	}
)
