package cmd

import (
	"context"

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
		Action: func(ctx context.Context, _ *cli.Command) error {
			err := user.Self(ctx)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
	}
)
