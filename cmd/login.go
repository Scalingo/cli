package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/session"
)

var (
	LoginCommand = cli.Command{
		Name:     "login",
		Category: "Global",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "api-key", Usage: "Authenticate with an API key instead of login/password"},
			cli.BoolFlag{Name: "ssh", Usage: "Login with you SSH identity instead of login/password"},
			cli.StringFlag{Name: "ssh-identity", Value: "ssh-agent", Usage: "Use a custom SSH key, only compatible if --ssh is set"},
		},
		Usage: "Login to Scalingo platform",
		Description: `
   Example
     'scalingo login'`,
		Action: func(c *cli.Context) {
			err := session.Login(session.LoginOpts{
				ApiKey:      c.String("api-key"),
				Ssh:         c.Bool("ssh"),
				SshIdentity: c.String("ssh-identity"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "login")
		},
	}
)
