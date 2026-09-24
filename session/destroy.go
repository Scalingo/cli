package session

import (
	"context"
	"net/http"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	scalingohttp "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

func DestroyToken(ctx context.Context) error {
	authenticator := &config.CliAuthenticator{}

	// A failure to revoke the token must not prevent the user from logging out locally.
	err := revokeToken(ctx, authenticator)
	if err != nil {
		config.C.Logger.Printf("Fail to revoke the API token: %+v\n", err)
		io.Warning("Fail to revoke the API token on Scalingo, it may still be active.")
		io.Warning("You can revoke it manually from the API tokens section of your account on the Scalingo dashboard.")
	}

	err = authenticator.RemoveAuth(ctx)
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

// revokeToken deletes the API token on the Auth API. Only the tokens created by the CLI at login have
// their ID stored locally: tokens provided by the user with `--api-token` and tokens stored by older
// versions of the CLI are not revoked.
func revokeToken(ctx context.Context, authenticator *config.CliAuthenticator) error {
	_, token, err := authenticator.LoadAuth(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "load authentication credentials")
	}
	if token == nil || token.ID == "" {
		return nil
	}

	c, err := config.ScalingoAuthClientFromToken(ctx, token.Token)
	if err != nil {
		return errors.Wrap(ctx, err, "create Scalingo auth client")
	}

	err = c.TokenDelete(ctx, token.ID)
	if err != nil {
		var requestFailedErr *scalingohttp.RequestFailedError
		// The token has already been revoked
		if errors.As(err, &requestFailedErr) &&
			(requestFailedErr.Code == http.StatusNotFound || requestFailedErr.Code == http.StatusUnauthorized) {
			return nil
		}
		return errors.Wrap(ctx, err, "delete API token")
	}

	return nil
}
