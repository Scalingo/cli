package user

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Self(ctx context.Context) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get scalingo API client")
	}
	user, err := c.Self(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get user info")
	}

	io.Statusf("You are logged in as %s (%s)\n", user.Username, user.Email)
	return nil
}
