package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	PsCommand = cli.Command{
		Name:     "ps",
		Category: "App Management",
		Usage:    "Display your application running processes",
		Flags:    []cli.Flag{appFlag},
		Description: `Display your application processes
	Example
	  'scalingo --app my-app ps'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "ps")
			} else if err := apps.Ps(currentApp); err != nil {
				errorQuit(err)
			}
		},
	}
)
