package cmd

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/codegangsta-cli"
)

func AuthenticateHook(c *cli.Context) error {
	apiKey := os.Getenv("SCALINGO_API_TOKEN")

	if apiKey == "" && config.AuthenticatedUser != nil {
		return nil
	}
	return session.Login(session.LoginOpts{ApiKey: apiKey})
}
