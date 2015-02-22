package cmd

import (
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/codegangsta-cli"
)

var (
	SignUpCommand = cli.Command{
		Name:     "signup",
		Category: "Global",
		Usage:    "Create your account of Scalingo",
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
