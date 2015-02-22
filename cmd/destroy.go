package cmd

import (
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	DestroyCommand = cli.Command{
		Name:        "destroy",
		Category:    "Global",
		ShortName:   "d",
		Usage:       "Destroy an app /!\\",
		Description: "Destroy an app /!\\ It is not reversible\n  Example:\n    'scalingo destroy my-app'",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "destroy")
			} else {
				err := apps.Destroy(c.Args()[0])
				if err != nil {
					errorQuit(err)
				}
			}
		},
	}
)
