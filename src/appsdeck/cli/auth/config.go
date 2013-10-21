package auth

import (
	"appsdeck/cli/api"
	"appsdeck/cli/constants"
	"encoding/json"
	"fmt"
	"os"
)

var (
	Config *AuthConfig
)

func init() {
	Config, err := LoadAuth()
	if err != nil {
		fmt.Println("You need to be authenticated to user Appsdeck client.\nNo account ? → https://appsdeck.eu/users/sign_up")
		Config, err = Authenticate()
		if err != nil {
			fmt.Println("An error occured :", err)
			os.Exit(1)
		} else {
			fmt.Printf("Hello %s %s, nice to see you !\n\n", Config.FirstName, Config.LastName)
		}
	}
	api.AuthToken = Config.AuthToken
	api.AuthEmail = Config.Email
}

type AuthConfig struct {
	AuthToken string `json:"authentication_token"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

func StoreAuth(authConfig *AuthConfig) error {
	// Check ~/.config/appsdeck
	if _, err := os.Stat(constants.ConfigDir); err != nil {
		if err, ok := err.(*os.PathError); ok {
			if err := os.MkdirAll(constants.ConfigDir, 0755); err != nil {
				return err
			}
		} else {
			fmt.Errorf("Error reaching config directory : %s", err)
		}
	}

	file, err := os.OpenFile(constants.AuthConfigFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(authConfig); err != nil {
		return err
	}
	return nil
}

func LoadAuth() (*AuthConfig, error) {
	file, err := os.OpenFile(constants.AuthConfigFile, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var authConfig AuthConfig
	dec := json.NewDecoder(file)
	if err := dec.Decode(&authConfig); err != nil {
		return nil, err
	}

	return &authConfig, nil
}
