package config

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/cli/config/auth"
	"github.com/Scalingo/go-scalingo/v6"
)

func TestStoreAuth(t *testing.T) {
	ctx := context.Background()
	u := &scalingo.User{
		Email:    "test@example.com",
		Username: "test",
	}
	authenticator := &CliAuthenticator{}

	// First creation
	err := authenticator.StoreAuth(ctx, u, "0123456789")
	require.NoError(t, err)
	clean()

	// Rewrite over an existing file
	err = authenticator.StoreAuth(ctx, u, "0123456789")
	require.NoError(t, err)
	err = authenticator.StoreAuth(ctx, u, "0123456789")
	require.NoError(t, err)
	clean()

	// Add an additional auth url
	err = authenticator.StoreAuth(ctx, u, "0123456789")
	require.NoError(t, err)
	C.ScalingoAuthURL = "api.scalingo2.dev"
	err = authenticator.StoreAuth(ctx, u, "0123456789")
	require.NoError(t, err)
	clean()
}

func TestExistingAuth(t *testing.T) {
	ctx := context.Background()
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
	err = authenticator.StoreAuth(ctx, u, "0123456789")
	require.NoError(t, err)

	currentAuth, err = existingAuth(ctx)
	json.Unmarshal(currentAuth.AuthConfigPerHost, &configPerHost)
	require.NoError(t, err)
	assert.Len(t, configPerHost, 1)
	assert.False(t, currentAuth.LastUpdate.IsZero())

	clean()
}
