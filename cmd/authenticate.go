package cmd

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/codegangsta-cli"
)

func AuthenticateHook(c *cli.Context) error {
	if config.AuthenticatedUser != nil {
		return nil
	}
	return session.Login(session.LoginOpts{})
}
