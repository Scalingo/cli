package cmd

import (
	"fmt"

	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/codegangsta-cli"
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
	}
)
