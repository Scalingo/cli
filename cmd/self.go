package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/user"
)

var (
	selfCommand = cli.Command{
		Name:        "self",
		Aliases: 	 []string{"whoami"},
		Category:    "Global",
		Usage:       "Get the logged in profile",
		Description: "Returns the logged in profile and print its username. Comes in handy when owning multiple accounts.",
		Action: func(c *cli.Context) {
			err := user.Self()
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
