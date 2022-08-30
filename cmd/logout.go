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
		Description: "Destroy login information stored on your computer",
		Action: func(c *cli.Context) error {
			currentUser, err := config.C.CurrentUser()
			if err != nil {
				errorQuit(err)
			}
			if currentUser == nil {
				io.Status("You are already logged out.")
				return nil
			}
			if err := session.DestroyToken(); err != nil {
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
