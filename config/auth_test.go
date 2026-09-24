package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/cli/config/auth"
	"github.com/Scalingo/go-scalingo/v11"
)

func TestStoreAuth(t *testing.T) {
	ctx := t.Context()
	u := &scalingo.User{
		Email:    "test@example.com",
		Username: "test",
	}
	authenticator := &CliAuthenticator{}

	token := auth.UserToken{Token: "0123456789"}

	// First creation
	err := authenticator.StoreAuth(ctx, u, token)
	require.NoError(t, err)
	clean()

	// Rewrite over an existing file
	err = authenticator.StoreAuth(ctx, u, token)
	require.NoError(t, err)
	err = authenticator.StoreAuth(ctx, u, token)
	require.NoError(t, err)
	clean()

	// Add an additional auth url
	err = authenticator.StoreAuth(ctx, u, token)
	require.NoError(t, err)
	C.ScalingoAuthURL = "api.scalingo2.dev"
	err = authenticator.StoreAuth(ctx, u, token)
	require.NoError(t, err)
	clean()
}

func TestExistingAuth(t *testing.T) {
	ctx := t.Context()
	u := &scalingo.User{
		Email:    "test@example.com",
		Username: "test",
	}
	authenticator := &CliAuthenticator{}

	// Before any auth
	currentAuth, err := existingAuth(ctx)
	require.NoError(t, err)

	var configPerHost auth.ConfigPerHostV2
	json.Unmarshal(currentAuth.AuthConfigPerHost, &configPerHost)
	assert.Empty(t, configPerHost)
	assert.True(t, currentAuth.LastUpdate.IsZero())

	// After one auth
	err = authenticator.StoreAuth(ctx, u, auth.UserToken{Token: "0123456789"})
	require.NoError(t, err)

	currentAuth, err = existingAuth(ctx)
	json.Unmarshal(currentAuth.AuthConfigPerHost, &configPerHost)
	require.NoError(t, err)
	assert.Len(t, configPerHost, 1)
	assert.False(t, currentAuth.LastUpdate.IsZero())

	clean()
}

func TestLoadAuth_TokenID(t *testing.T) {
	ctx := t.Context()
	u := &scalingo.User{
		Email:    "test@example.com",
		Username: "test",
	}
	authenticator := &CliAuthenticator{}
	defer clean()

	t.Run("it should load the token ID stored at login", func(t *testing.T) {
		err := authenticator.StoreAuth(ctx, u, auth.UserToken{Token: "0123456789", ID: "token-id"})
		require.NoError(t, err)

		_, token, err := authenticator.LoadAuth(ctx)
		require.NoError(t, err)
		assert.Equal(t, "0123456789", token.Token)
		assert.Equal(t, "token-id", token.ID)
	})

	t.Run("it should load a token without ID", func(t *testing.T) {
		err := authenticator.StoreAuth(ctx, u, auth.UserToken{Token: "0123456789"})
		require.NoError(t, err)

		_, token, err := authenticator.LoadAuth(ctx)
		require.NoError(t, err)
		assert.Equal(t, "0123456789", token.Token)
		assert.Empty(t, token.ID)
	})
}
