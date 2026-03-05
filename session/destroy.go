package session

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func DestroyToken(ctx context.Context) error {
	authenticator := &config.CliAuthenticator{}
	err := authenticator.RemoveAuth(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "remove local authentication credentials")
	}
	return nil
}
