package session

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func DestroyToken(ctx context.Context) error {
	authenticator := &config.CliAuthenticator{}
	err := authenticator.RemoveAuth(ctx)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
