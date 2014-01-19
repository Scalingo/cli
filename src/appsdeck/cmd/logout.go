package cmd

import (
	"appsdeck/session"
	"fmt"
	"github.com/codegangsta/cli"
)

var (
	LogoutCommand = cli.Command{
		Name:        "logout",
		Usage:       "Logout from Appsdeck",
		Description: "Destroy login information stored on your computer",
		Action: func(c *cli.Context) {
			if err := session.DestroyToken(); err != nil {
				panic(err)
			}
			fmt.Println("Appsdeck credentials have been deleted.")
		},
	}
)
