package config

import (
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/users"
	"gopkg.in/errgo.v1"
)

func StoreAuth(user *users.User) error {
	// Check ~/.config/scalingo
	if _, err := os.Stat(C.ConfigDir); err != nil {
		if err, ok := err.(*os.PathError); ok {
			if err := os.MkdirAll(C.ConfigDir, 0755); err != nil {
				return errgo.Mask(err, errgo.Any)
			}
		} else {
			return errgo.Notef(err, "error reaching config directory")
		}
	}

	file, err := os.OpenFile(C.AuthFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(user); err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}

func LoadAuth() (*users.User, error) {
	file, err := os.OpenFile(C.AuthFile, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	var user *users.User
	dec := json.NewDecoder(file)
	if err := dec.Decode(&user); err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return user, nil
}
