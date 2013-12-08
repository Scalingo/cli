package cmd

import (
	"appsdeck/apps"
	"appsdeck/auth"
	"github.com/codegangsta/cli"
)

var (
	AppsCommand = cli.Command{
		Name:      "apps",
		ShortName: "a",
		Usage:     "Manage your apps",
		Flags: []cli.Flag{
			cli.BoolFlag{"list", "List your apps"},
		},
		Action: func(c *cli.Context) {
			auth.InitAuth()
			if c.Bool("list") {
				if err := apps.List(); err != nil {
					errorQuit(err)
				}
			} else {
				cli.ShowCommandHelp(c, "apps")
			}
		},
	}
)
