package config

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Scalingo/envconfig"
	"github.com/Scalingo/go-scalingo"
	"github.com/stvp/rollbar"
	"gopkg.in/errgo.v1"
)

type ConfigFile struct {
	Region string `json:"region"`
}

type Config struct {
	ApiVersion           string
	DisableInteractive   bool
	DisableUpdateChecker bool
	UnsecureSsl          bool
	RollbarToken         string

	// Override region configuration
	ScalingoApiUrl  string
	ScalingoAuthUrl string
	ScalingoDbUrl   string
	ScalingoRegion  string
	ScalingoSshHost string

	// Configuration files
	ConfigDir      string
	AuthFile       string
	LogFile        string
	ConfigFilePath string
	ConfigFile     ConfigFile

	// Cache related files
	CacheDir         string
	RegionsCachePath string

	// Logging
	logFile *os.File
	Logger  *log.Logger
}

var (
	env = map[string]string{
		"SCALINGO_AUTH_URL":  "https://auth.scalingo.com",
		"SCALINGO_API_URL":   "",
		"SCALINGO_DB_URL":    "",
		"SCALINGO_SSH_HOST":  "scalingo.com:22",
		"SCALINGO_REGION":    "",
		"API_VERSION":        "1",
		"UNSECURE_SSL":       "false",
		"ROLLBAR_TOKEN":      "",
		"CONFIG_DIR":         ".config/scalingo",
		"CACHE_DIR":          ".cache/scalingo",
		"AUTH_FILE":          "auth",
		"CONFIG_FILE_PATH":   "config.json",
		"REGIONS_CACHE_PATH": "regions.json",
		"LOG_FILE":           "local.log",
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
	env["CONFIG_FILE_PATH"] = filepath.Join(env["CONFIG_DIR"], env["CONFIG_FILE_PATH"])
	env["LOG_FILE"] = filepath.Join(env["CONFIG_DIR"], env["LOG_FILE"])

	env["CACHE_DIR"] = filepath.Join(home, env["CACHE_DIR"])
	env["REGIONS_CACHE_PATH"] = filepath.Join(env["CACHE_DIR"], env["REGIONS_CACHE_PATH"])

	for k := range env {
		vEnv := os.Getenv(k)
		if vEnv == "" {
			os.Setenv(k, env[k])
		}
	}

	err := envconfig.Process("", &C)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to initialize configuration: %v\n", err)
		os.Exit(1)
	}

	err = os.MkdirAll(C.ConfigDir, 0750)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to create configuration directory: %v\n", err)
		os.Exit(1)
	}
	err = os.MkdirAll(C.CacheDir, 0750)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to create configuration directory: %v\n", err)
		os.Exit(1)
	}

	C.logFile, err = os.OpenFile(C.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to open log file: %s, disabling logging.\n", C.LogFile)
	}
	C.Logger = log.New(C.logFile, "", log.LstdFlags)

	rollbar.Token = C.RollbarToken
	rollbar.Platform = "client"
	rollbar.Environment = "production"
	rollbar.ErrorWriter = C.logFile

	TlsConfig = &tls.Config{}
	if C.UnsecureSsl {
		TlsConfig.InsecureSkipVerify = true
		TlsConfig.MinVersion = tls.VersionTLS10
	}

	// Read region from the configuration file
	fd, err := os.Open(C.ConfigFilePath)
	if err == nil {
		json.NewDecoder(fd).Decode(&C.ConfigFile)
	}
	if C.ScalingoRegion == "" {
		C.ScalingoRegion = C.ConfigFile.Region
	}
}

func (config Config) CurrentUser() (*scalingo.User, error) {
	authenticator := &CliAuthenticator{}
	user, _, err := authenticator.LoadAuth()
	return user, err
}

type ClientConfigOpts struct {
	Region   string
	APIToken string
	AuthOnly bool
}

func (config Config) scalingoClientBaseConfig(opts ClientConfigOpts) scalingo.ClientConfig {
	return scalingo.ClientConfig{
		TLSConfig:    TlsConfig,
		Region:       opts.Region,
		AuthEndpoint: config.ScalingoAuthUrl,
		APIToken:     opts.APIToken,
		UserAgent:    "Scalingo CLI v" + Version,
	}
}

func (config Config) scalingoClientConfig(opts ClientConfigOpts) (scalingo.ClientConfig, error) {
	c := config.scalingoClientBaseConfig(opts)
	if !opts.AuthOnly {
		if config.ScalingoApiUrl != "" && config.ScalingoDbUrl != "" {
			c.APIEndpoint = config.ScalingoApiUrl
			c.DatabaseAPIEndpoint = config.ScalingoDbUrl
		} else {
			region, err := GetRegion(config, config.ScalingoRegion, GetRegionOpts{
				Token: opts.APIToken,
			})
			if err != nil {
				return c, errgo.Notef(err, "fail to get region '%v' specifications", config.ScalingoRegion)
			}
			c.APIEndpoint = region.API
			c.DatabaseAPIEndpoint = region.DatabaseAPI
		}
	}
	return c, nil
}

func ScalingoClientFromToken(token string) (*scalingo.Client, error) {
	config, err := C.scalingoClientConfig(ClientConfigOpts{APIToken: token})
	if err != nil {
		return nil, errgo.Notef(err, "fail to create Scalingo client")
	}
	return scalingo.New(config)
}

func ScalingoAuthClientFromToken(token string) (*scalingo.Client, error) {
	config := C.scalingoClientBaseConfig(ClientConfigOpts{
		AuthOnly: true, APIToken: token,
	})
	return scalingo.New(config)
}

func ScalingoAuthClient() (*scalingo.Client, error) {
	auth := &CliAuthenticator{}
	_, token, err := auth.LoadAuth()
	if err != nil {
		return nil, errgo.Notef(err, "fail to load authentication")
	}
	return ScalingoAuthClientFromToken(token.Token)
}

func ScalingoClient() (*scalingo.Client, error) {
	authenticator := &CliAuthenticator{}
	_, token, err := authenticator.LoadAuth()
	if err != nil {
		return nil, errgo.Notef(err, "fail to load credentials")
	}
	config, err := C.scalingoClientConfig(ClientConfigOpts{APIToken: token.Token})
	if err != nil {
		return nil, errgo.Notef(err, "fail to create Scalingo client")
	}
	return scalingo.New(config)
}

func ScalingoClientForRegion(region string) (*scalingo.Client, error) {
	authenticator := &CliAuthenticator{}
	_, token, err := authenticator.LoadAuth()
	if err != nil {
		return nil, errgo.Notef(err, "fail to load credentials")
	}

	config := C.scalingoClientBaseConfig(ClientConfigOpts{
		Region:   region,
		APIToken: token.Token,
	})

	return scalingo.New(config)
}

func ScalingoUnauthenticatedAuthClient() (*scalingo.Client, error) {
	config := C.scalingoClientBaseConfig(ClientConfigOpts{AuthOnly: true})
	return scalingo.New(config)
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
