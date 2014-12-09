package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Scalingo/envconfig"
	"github.com/stvp/rollbar"
)

type Config struct {
	ApiUrl       string
	ApiPrefix    string
	SshHost      string
	UnsecureSsl  bool
	RollbarToken string
	ConfigDir    string
	AuthFile     string
	LogFile      string
}

var (
	env = map[string]string{
		"API_URL":       "https://scalingo-api-production.appsdeck.eu",
		"SSH_HOST":      "appsdeck.eu:22",
		"API_PREFIX":    "/v1",
		"UNSECURE_SSL":  "false",
		"ROLLBAR_TOKEN": "",
		"CONFIG_DIR":    ".config/scalingo",
		"AUTH_FILE":     "auth",
		"LOG_FILE":      "local.log",
	}
	C         Config
	TlsConfig *tls.Config
)

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		panic("The HOME environment variable must be defined")
	}

	env["CONFIG_DIR"] = filepath.Join(home, env["CONFIG_DIR"])
	env["AUTH_FILE"] = filepath.Join(env["CONFIG_DIR"], env["AUTH_FILE"])
	env["LOG_FILE"] = filepath.Join(env["CONFIG_DIR"], env["LOG_FILE"])

	for k := range env {
		vEnv := os.Getenv(k)
		if vEnv == "" {
			os.Setenv(k, env[k])
		}
	}

	envconfig.Process("", &C)

	logFile, err := os.OpenFile(C.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to open log file: %s, disabling logging.\n", C.LogFile)
	}

	rollbar.Token = C.RollbarToken
	rollbar.Platform = "client"
	rollbar.Environment = "production"
	rollbar.ErrorWriter = logFile

	TlsConfig = &tls.Config{}
	if C.UnsecureSsl {
		TlsConfig.InsecureSkipVerify = true
		TlsConfig.MinVersion = tls.VersionTLS10
	} else {
		certChain := decodePem(x509Chain)
		TlsConfig.RootCAs = x509.NewCertPool()
		for _, cert := range certChain.Certificate {
			x509Cert, err := x509.ParseCertificate(cert)
			if err != nil {
				panic(err)
			}
			TlsConfig.RootCAs.AddCert(x509Cert)
		}
		TlsConfig.BuildNameToCertificate()
	}
}

func GenTLSConfig(serverName string) *tls.Config {
	tlsConfig := &tls.Config{}
	*tlsConfig = *TlsConfig
	tlsConfig.ServerName = serverName
	return tlsConfig
}
