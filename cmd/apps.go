package cmd

import (
	"github.com/Appsdeck/appsdeck/apps"
	"github.com/Appsdeck/appsdeck/auth"
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
