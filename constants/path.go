package constants

import (
	"os"
	"path"
)

var (
	ConfigDir      = ".config/appsdeck"
	AuthConfigFile = ".config/appsdeck/auth"
)

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		panic("The HOME environment variable must be defined")
	}

	ConfigDir = path.Join(home, ConfigDir)
	AuthConfigFile = path.Join(home, AuthConfigFile)
}
