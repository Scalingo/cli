package cmd

import (
	"github.com/Scalingo/cli/addons"
	"github.com/codegangsta/cli"
)

var (
	AddonsListCommand = cli.Command{
		Name:        "addons-list",
		Category:    "Addons - Global",
		Description: "List all addons you can add to your app.",
		Usage:       "List all addons",
		Action: func(c *cli.Context) {
			if err := addons.List(); err != nil {
				errorQuit(err)
			}
		},
	}
)
