package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/user"
)

var (
	selfCommand = cli.Command{
		Name:        "self",
		Aliases:     []string{"whoami"},
		Category:    "Global",
		Usage:       "Get the logged in profile",
		Description: "Returns the logged in profile and print its username. Comes in handy when owning multiple accounts.",
		Action: func(ctx context.Context, c *cli.Command) error {
			err := user.Self(c.Context)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
	}
)
