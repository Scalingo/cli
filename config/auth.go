package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Scalingo/cli/constants"
	"github.com/Scalingo/cli/users"
)

func StoreAuth(user *users.User) error {
	// Check ~/.config/scalingo
	if _, err := os.Stat(constants.ConfigDir); err != nil {
		if err, ok := err.(*os.PathError); ok {
			if err := os.MkdirAll(constants.ConfigDir, 0755); err != nil {
				return err
			}
		} else {
			fmt.Errorf("Error reaching config directory: %s", err)
		}
	}

	file, err := os.OpenFile(constants.AuthConfigFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(user); err != nil {
		return err
	}
	return nil
}

func LoadAuth() (*users.User, error) {
	file, err := os.OpenFile(constants.AuthConfigFile, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var user *users.User
	dec := json.NewDecoder(file)
	if err := dec.Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}
