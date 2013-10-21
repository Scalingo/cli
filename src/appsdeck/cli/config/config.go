package config

import (
	"crypto/tls"
	"os"
)

var (
	TlsConfig *tls.Config
	C map[string]string

	defaultValues = map[string]string{
		"APPSDECK_LOG": "https://logs.appsdeck.eu",
		"APPSDECK_API": "https://appsdeck.eu",
		"UNSECURE_SSL": "false",
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

	TlsConfig = &tls.Config{}
	if C["UNSECURE_SSL"] == "true" {
		TlsConfig.InsecureSkipVerify = true
	}
}
