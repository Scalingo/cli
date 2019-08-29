package config

import (
	"os"
	"path/filepath"
)

var (
	testConfig = Config{
		ScalingoApiUrl: "api.scalingo.dev",
		AuthFile:       filepath.Join(os.TempDir(), "test-scalingo-auth"),
	}
)

func init() {
	C = testConfig
}

func clean() {
	os.Remove(testConfig.AuthFile)
}
