package cmd

import (
	"appsdeck/cli/apps"
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
