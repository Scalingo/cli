package cmd

import (
	"appsdeck/cli/apps"
	"github.com/Appsdeck/cli"
)

var (
	AppsCommand = cli.Command{
		Name:	"apps",
		ShortName: "a",
		Usage: "Manage your apps",
		Flags: []cli.Flag{
			cli.BoolFlag{"list", "List your apps"},
		},
		Action: func(c *cli.Context) {
			if c.Bool("list") {
				apps.List()
			}
		},
	}
)

