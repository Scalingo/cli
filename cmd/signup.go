package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/session"
)

var (
	SignUpCommand = cli.Command{
		Name:     "signup",
		Category: "Global",
		Usage:    "Create your Scalingo account",
		Description: `
   Example
     'scalingo signup'`,
		Action: func(c *cli.Context) {
			err := session.SignUp()
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
