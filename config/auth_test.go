package config

import (
	"testing"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
)

var (
	u = &scalingo.User{
		Email:     "test@example.com",
		Username:  "test",
		AuthToken: "0123456789",
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
	auth, err := existingAuth()
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	if len(auth.AuthConfigPerHost) > 0 {
		t.Errorf("want auth.AuthConfigPerHost = [], got %v", auth.AuthConfigPerHost)
	}
	if !auth.LastUpdate.IsZero() {
		t.Errorf("auth should never have been updated: %v", auth.LastUpdate)
	}

	// After one auth
	err = Authenticator.StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	auth, err = existingAuth()
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	if len(auth.AuthConfigPerHost) != 1 {
		t.Errorf("want len(auth.AuthConfigPerHost) = 1, got %v", auth.AuthConfigPerHost)
	}
	if auth.LastUpdate.IsZero() {
		t.Errorf("auth should have been updated: %v", auth.LastUpdate)
	}

	clean()
}
