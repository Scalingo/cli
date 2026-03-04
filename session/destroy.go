package session

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func DestroyToken(ctx context.Context) error {
	authenticator := &config.CliAuthenticator{}
	err := authenticator.RemoveAuth(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}
	return nil
}
