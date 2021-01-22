package config

import (
	"encoding/json"
	"testing"

	"github.com/Scalingo/cli/config/auth"
	"github.com/Scalingo/go-scalingo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAuth(t *testing.T) {
	u := &scalingo.User{
		Email:    "test@example.com",
		Username: "test",
	}
	authenticator := &CliAuthenticator{}

	// First creation
	err := authenticator.StoreAuth(u, "0123456789")
	require.NoError(t, err)
	clean()

	// Rewrite over an existing file
	err = authenticator.StoreAuth(u, "0123456789")
	require.NoError(t, err)
	err = authenticator.StoreAuth(u, "0123456789")
	require.NoError(t, err)
	clean()

	// Add an additional auth url
	err = authenticator.StoreAuth(u, "0123456789")
	require.NoError(t, err)
	C.ScalingoAuthUrl = "api.scalingo2.dev"
	err = authenticator.StoreAuth(u, "0123456789")
	require.NoError(t, err)
	clean()
}

func TestExistingAuth(t *testing.T) {
	u := &scalingo.User{
		Email:    "test@example.com",
		Username: "test",
	}
	authenticator := &CliAuthenticator{}

	// Before any auth
	currentAuth, err := existingAuth()
	require.NoError(t, err)

	var configPerHost auth.ConfigPerHostV2
	json.Unmarshal(currentAuth.AuthConfigPerHost, &configPerHost)
	assert.Empty(t, configPerHost)
	assert.True(t, currentAuth.LastUpdate.IsZero())

	// After one auth
	err = authenticator.StoreAuth(u, "0123456789")
	require.NoError(t, err)

	currentAuth, err = existingAuth()
	json.Unmarshal(currentAuth.AuthConfigPerHost, &configPerHost)
	require.NoError(t, err)
	assert.Len(t, configPerHost, 1)
	assert.False(t, currentAuth.LastUpdate.IsZero())

	clean()
}
