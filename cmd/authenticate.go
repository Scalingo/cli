package cmd

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/session"
	"github.com/urfave/cli"
)

func AuthenticateHook(c *cli.Context) error {
	token := os.Getenv("SCALINGO_API_TOKEN")

	if token == "" && config.AuthenticatedUser != nil {
		return nil
	}
	return session.Login(session.LoginOpts{APIToken: token})
}
