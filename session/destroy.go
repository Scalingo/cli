package session

import (
	"os"

	"github.com/Scalingo/cli/config"
)

func DestroyToken() error {
	if _, err := os.Stat(config.C.AuthFile); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(config.C.AuthFile)
}
