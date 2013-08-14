package config

import (
	"os"
)

var (
	C map[string]string

	defaultValues = map[string]string{
		"APPSDECK_LOG": "logs.appsdeck.eu",
		"APPSDECK_API": "appsdeck.eu",
	}
)

func init() {
	C = make(map[string]string)
	for varName, defaultValue := range defaultValues {
		C[varName] = os.Getenv(varName)
		if C[varName] == "" {
			C[varName] = defaultValue
		}
	}
}
