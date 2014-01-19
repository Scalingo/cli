package cmd

import (
	"appsdeck/apps"
	"appsdeck/auth"
	"github.com/codegangsta/cli"
)

var (
	AppsCommand = cli.Command{
		Name:        "apps",
		ShortName:   "a",
		Description: "List your apps and give some details about them",
		Usage:       "List your apps",
		Action: func(c *cli.Context) {
			auth.InitAuth()
			if err := apps.List(); err != nil {
				errorQuit(err)
			}
		},
	}
)
