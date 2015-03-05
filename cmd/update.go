package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/update"
)

var (
	UpdateCommand = cli.Command{
		Name:     "update",
		Category: "CLI Internals",
		Usage:    "Update 'scalingo' client",
		Description: `Update 'scalingo' client
   Example
     'scalingo update'`,
		Action: func(c *cli.Context) {
			err := update.Check()
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
