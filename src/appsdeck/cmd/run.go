package cmd

import (
	"appsdeck/appdetect"
	"appsdeck/apps"
	"github.com/codegangsta/cli"
)

var (
	RunCommand = cli.Command{
		Name:      "run",
		ShortName: "r",
		Usage:     "Run any command for your app",
		Description: `Run command in current app context, your application
   environment will be loaded and you can execute any task (example
   'rake' or any database-related task)`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			if len(c.Args()) == 0 {
				cli.ShowCommandHelp(c, "run")
			} else if err := apps.Run(currentApp, c.Args()); err != nil {
				errorQuit(err)
			}
		},
	}
)
