package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/session"
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
