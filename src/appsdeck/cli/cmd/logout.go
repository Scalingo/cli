package cmd

import (
	"appsdeck/cli/auth"
	"fmt"
	"github.com/codegangsta/cli"
)

var (
	LogoutCommand = cli.Command{
		Name:  "logout",
		Usage: "Logout from Appsdeck",
		Action: func(c *cli.Context) {
			if err := auth.DestroyToken(); err != nil {
				panic(err)
			}
			fmt.Println("Appsdeck credentials have been deleted.")
		},
	}
)
