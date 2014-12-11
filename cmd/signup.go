package cmd

import (
	"github.com/Scalingo/cli/session"
	"github.com/codegangsta/cli"
)

var (
	SignUpCommand = cli.Command{
		Name:  "signup",
		Usage: "Create your account of Scalingo",
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
