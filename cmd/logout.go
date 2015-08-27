package cmd

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/session"
)

var (
	LogoutCommand = cli.Command{
		Name:        "logout",
		Category:    "Global",
		Usage:       "Logout from Scalingo",
		Description: "Destroy login information stored on your computer",
		Action: func(c *cli.Context) {
			if err := session.DestroyToken(); err != nil {
				panic(err)
			}
			fmt.Println("Scalingo credentials have been deleted.")
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "logout")
		},
	}
)
