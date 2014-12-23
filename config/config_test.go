package config

import "os"

var (
	testConfig = Config{
		apiHost:  "scalingo.dev",
		AuthFile: "/tmp/test-scalingo-auth",
	}
)

func init() {
	C = testConfig
}

func clean() {
	os.Remove(testConfig.AuthFile)
}
