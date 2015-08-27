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
		Description: "Create a new app:\n   Example:\n     'scalingo create mynewapp'\n     'scalingo create mynewapp --remote \"staging\"'",
		Flags:       []cli.Flag{cli.StringFlag{Name: "remote", Value: "scalingo", Usage: "Remote to add to your current git repository", EnvVar: ""}},
		Usage:       "Create a new app",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "create")
			} else {
				err := apps.Create(c.Args()[0], c.String("remote"))
				if err != nil {
					errorQuit(err)
				}
			}
		},
	}
)
