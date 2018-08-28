package cmd

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/urfave/cli"
)

var (
	selfCommand = cli.Command{
		Name:        "self",
		Category:    "Global",
		Usage:       "Get the logged in profile",
		Description: "Returns the logged in profile and print its username. Comes in handy when owning multiple accounts.",
		Before:      AuthenticateHook,
		Action: func(c *cli.Context) {
			if config.AuthenticatedUser != nil {
				io.Statusf("You are logged in as %s (%s)\n", config.AuthenticatedUser.Username, config.AuthenticatedUser.Email)
				return
			}
			io.Status("Currently unauthenticated")
		},
	}
	whoamiCommand = cli.Command{
		Name:        "whoami",
		Category:    selfCommand.Category,
		Usage:       selfCommand.Usage,
		Description: selfCommand.Description,
		Before:      selfCommand.Before,
		Action:      selfCommand.Action,
	}
)
