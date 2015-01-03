package cmd

import (
	"github.com/Scalingo/cli/session"
	"github.com/codegangsta/cli"
)

var (
	LoginCommand = cli.Command{
		Name:     "login",
		Category: "Global",
		Usage:    "Login to Scalingo platform",
		Description: `
   Example
     'scalingo login'`,
		Action: func(c *cli.Context) {
			err := session.Login()
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
