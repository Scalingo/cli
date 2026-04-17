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

	// We want to delete the regions cache so that the cache is created again during the next login.
	// This is important for people with two different Scalingo accounts, one with and one without access to the osc-secnum-fr1 region.
	// If we don't, there is a risk that when the client login with the osc-secnum-fr1 account, the regions cache does not contain the osc-secnum-fr1 region and client cannot contact this region.
	// Ref. https://github.com/Scalingo/cli/issues/1057
	err = config.DeleteRegionsCache(ctx, config.C)
	if err != nil {
		return errors.Wrap(ctx, err, "remove local regions cache")
	}

	return nil
}
