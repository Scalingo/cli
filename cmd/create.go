package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/apps"
)

var (
	CreateCommand = cli.Command{
		Name:        "create",
		Category:    "Global",
		ShortName:   "c",
		Description: "Create a new app:\n   Example:\n     'scalingo create mynewapp'",
		Usage:       "Create a new app",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "create")
			} else {
				err := apps.Create(c.Args()[0])
				if err != nil {
					errorQuit(err)
				}
			}
		},
	}
)
