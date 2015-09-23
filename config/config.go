package config

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/envconfig"
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/stvp/rollbar"
	"github.com/Scalingo/go-scalingo"
)

type Config struct {
	ApiUrl       string
	apiHost      string
	ApiVersion   string
	SshHost      string
	UnsecureSsl  bool
	RollbarToken string
	ConfigDir    string
	AuthFile     string
	LogFile      string
	logFile      *os.File
	Logger       *log.Logger
}

var (
	env = map[string]string{
		"API_URL":       "https://api.scalingo.com",
		"SSH_HOST":      "appsdeck.eu:22",
		"API_VERSION":   "1",
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
	home := HomeDir()
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

	err := os.MkdirAll(C.ConfigDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to create configuration directory: %v\n", err)
		os.Exit(1)
	}

	C.logFile, err = os.OpenFile(C.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to open log file: %s, disabling logging.\n", C.LogFile)
	}
	C.Logger = log.New(C.logFile, "", log.LstdFlags)

	u, err := url.Parse(C.ApiUrl)
	if err != nil {
		panic("API_URL is not a valid URL " + err.Error())
	}

	C.apiHost = strings.Split(u.Host, ":")[0]

	rollbar.Token = C.RollbarToken
	rollbar.Platform = "client"
	rollbar.Environment = "production"
	rollbar.ErrorWriter = C.logFile

	TlsConfig = &tls.Config{}
	if C.UnsecureSsl {
		TlsConfig.InsecureSkipVerify = true
		TlsConfig.MinVersion = tls.VersionTLS10
	}

	scalingo.ApiAuthenticator = Authenticator
	scalingo.ApiUrl = C.ApiUrl
	scalingo.ApiVersion = C.ApiVersion
}

func HomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
