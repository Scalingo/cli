package user

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Self() error {
	c, err := config.ScalingoAuthClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo API client")
	}
	user, err := c.Self()
	if err != nil {
		return errgo.Notef(err, "fail to get user info")
	}

	io.Statusf("You are logged in as %s (%s)\n", user.Username, user.Email)
	return nil
}
