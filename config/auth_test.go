package config

import (
	"encoding/json"
	"testing"

	"github.com/Scalingo/cli/config/auth"
	"github.com/Scalingo/go-scalingo"
)

var (
	u = &scalingo.User{
		Email:               "test@example.com",
		Username:            "test",
		AuthenticationToken: "0123456789",
	}
)

func TestStoreAuth(t *testing.T) {
	// First creation
	err := Authenticator.StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	clean()

	// Rewrite over an existing file
	err = Authenticator.StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	err = Authenticator.StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	clean()

	// Add an additional api url
	err = Authenticator.StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	C.apiHost = "scalingo2.dev"
	err = Authenticator.StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	clean()
}

func TestExistingAuth(t *testing.T) {
	// Before any auth
	currentAuth, err := existingAuth()
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	var configPerHost auth.ConfigPerHostV1
	json.Unmarshal(currentAuth.AuthConfigPerHost, &configPerHost)

	if len(configPerHost) > 0 {
		t.Errorf("want auth.AuthConfigPerHost = [], got %v", configPerHost)
	}
	if !currentAuth.LastUpdate.IsZero() {
		t.Errorf("auth should never have been updated: %v", currentAuth.LastUpdate)
	}

	// After one auth
	err = Authenticator.StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}

	currentAuth, err = existingAuth()
	json.Unmarshal(currentAuth.AuthConfigPerHost, &configPerHost)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	if len(configPerHost) != 1 {
		t.Errorf("want len(auth.AuthConfigPerHost) = 1, got %v", configPerHost)
	}
	if currentAuth.LastUpdate.IsZero() {
		t.Errorf("auth should have been updated: %v", currentAuth.LastUpdate)
	}

	clean()
}
