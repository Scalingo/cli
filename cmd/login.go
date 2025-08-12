package cmd

import (
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/session"
)

var (
	LoginCommand = cli.Command{
		Name:     "login",
		Category: "Global",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "api-token", Usage: "Authenticate with a token instead of login/password or SSH"},
			&cli.BoolFlag{Name: "ssh", Usage: "Login with you SSH identity instead of login/password"},
			&cli.StringFlag{Name: "ssh-identity", Value: "ssh-agent", Usage: "Use a custom SSH key, only compatible if --ssh is set"},
			&cli.BoolFlag{Name: "password-only", Usage: "Login with login/password without testing SSH connection"},
		},
		Usage:       "Login to Scalingo platform",
		Description: "Login to Scalingo platform",
		Action: func(c *cli.Context) error {
			if c.Bool("ssh") && c.Bool("password-only") {
				errorQuit(c.Context, errors.New("you cannot use both --ssh and --password-only at the same time"))
			}

			err := session.Login(c.Context, session.LoginOpts{
				APIToken:     c.String("api-token"),
				PasswordOnly: c.Bool("password-only"),
				SSH:          c.Bool("ssh"),
				SSHIdentity:  c.String("ssh-identity"),
			})
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "login")
		},
	}
)
