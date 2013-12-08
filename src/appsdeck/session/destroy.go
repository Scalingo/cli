package session

import (
	"appsdeck/constants"
	"os"
)

func DestroyToken() error {
	return os.Remove(constants.AuthConfigFile)
}
