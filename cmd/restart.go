package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/codegangsta/cli"
)

var (
	RestartCommand = cli.Command{
		Name:  "restart",
		Usage: "Restart processes of your app",
		Flags: []cli.Flag{cli.BoolFlag{Name: "synchronous", Usage: "Do the restart synchronously", EnvVar: ""}},
		Description: `Restart one or several process or your application:
	Example
	  ## Restart all the processes
	  scalingo --app my-app restart
		## Restart all the web processes
	  scalingo --app my-app restart web
		## Restart a specific container
	  scalingo --app my-app restart web-1`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if err := apps.Restart(currentApp, c.Bool("synchronous"), c.Args()); err != nil {
				errorQuit(err)
			}
		},
	}
)
