package config

import (
	"reflect"
	"testing"
)

var (
	u = &users.User{
		Email:     "test@example.com",
		Username:  "test",
		AuthToken: "0123456789",
	}
)

func TestStoreAuth(t *testing.T) {
	// First creation
	err := StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	clean()

	// Rewrite over an existing file
	err = StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	err = StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	clean()

	// Add an additional api url
	err = StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	C.apiHost = "scalingo2.dev"
	err = StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	clean()
}

func TestLoadAuth(t *testing.T) {
	// Load without Store should return nil User
	u, err := LoadAuth()
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	if u != nil {
		t.Errorf("%v should be nil", u)
	}

	// Load after storage of credentials
	err = StoreAuth(u)
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	uLoad, err := LoadAuth()
	if err != nil {
		t.Errorf("%v should be nil", err)
	}
	if !reflect.DeepEqual(uLoad, u) {
		t.Errorf("want %v, got %v", u, uLoad)
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
	err = StoreAuth(u)
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
