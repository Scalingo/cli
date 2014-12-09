package cmd

import (
	"github.com/Scalingo/cli/update"
	"github.com/codegangsta/cli"
)

var (
	UpdateCommand = cli.Command{
		Name:  "update",
		Usage: "Update 'scalingo' client",
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
