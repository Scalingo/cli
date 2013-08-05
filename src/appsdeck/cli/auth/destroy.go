package auth

import (
	"appsdeck/cli/constants"
	"os"
)

func DestroyToken() error {
	return os.Remove(constants.AuthConfigFile)
}
