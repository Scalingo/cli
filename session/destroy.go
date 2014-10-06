package session

import (
	"github.com/Scalingo/cli/constants"
	"os"
)

func DestroyToken() error {
	if _, err := os.Stat(constants.AuthConfigFile); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(constants.AuthConfigFile)
}
