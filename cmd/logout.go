package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/session"
)

var (
	LogoutCommand = cli.Command{
		Name:        "logout",
		Category:    "Global",
		Usage:       "Logout from Scalingo",
		Description: "Remove login information stored on your computer",
		Action: func(c *cli.Context) error {
			ctx := c.Context
			currentUser, err := config.C.CurrentUser(ctx)
			if err != nil {
				errorQuit(ctx, err)
			}
			if currentUser == nil {
				io.Status("You are already logged out.")
				return nil
			}
			err = session.DestroyToken(ctx)
			if err != nil {
				panic(err)
			}
			io.Status("Scalingo credentials have been deleted.")
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logout")
		},
	}
)
