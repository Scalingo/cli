package cmd

import (
	"github.com/Scalingo/cli/user"
	"github.com/urfave/cli"
)

var (
	selfCommand = cli.Command{
		Name:        "self",
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

	// whoami is an alias of the command selfCommand
	whoamiCommand = cli.Command{
		Name:        "whoami",
		Category:    selfCommand.Category,
		Usage:       selfCommand.Usage,
		Description: selfCommand.Description,
		Before:      selfCommand.Before,
		Action:      selfCommand.Action,
	}
)
