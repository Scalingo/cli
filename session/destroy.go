package session

import (
	"github.com/Appsdeck/appsdeck/constants"
	"os"
)

func DestroyToken() error {
	if _, err := os.Stat(constants.AuthConfigFile); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(constants.AuthConfigFile)
}
