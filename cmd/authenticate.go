package cmd

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/session"
	"github.com/urfave/cli"
)

func AuthenticateHook(c *cli.Context) error {
	token := os.Getenv("SCALINGO_API_TOKEN")

	currentUser, err := config.C.CurrentUser()
	if err == nil && currentUser != nil {
		return nil
	}
	return session.Login(session.LoginOpts{APIToken: token})
}
