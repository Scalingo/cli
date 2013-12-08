package cmd

import (
	"appsdeck/apps"
	"appsdeck/auth"
	"github.com/codegangsta/cli"
)

var (
	CreateCommand = cli.Command{
		Name:        "create",
		ShortName:   "c",
		Description: "Create a new app",
		Usage:       "appsdeck create <name>",
		Action: func(c *cli.Context) {
			auth.InitAuth()
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "create")
			} else {
				apps.Create(c.Args()[0])
			}
		},
	}
)
