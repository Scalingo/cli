package config

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/stvp/rollbar"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

type ConfigFile struct {
	Region string `json:"region"`
}

var (
	// Injected from compiler when making a release
	RollbarToken = ""
)

type Config struct {
	APIVersion           string `envconfig:"API_VERSION"`
	DisableInteractive   bool   `envconfig:"DISABLE_INTERACTIVE"`
	DisableUpdateChecker bool   `envconfig:"DISABLE_UPDATE_CHECKER"`
	UnsecureSsl          bool   `envconfig:"UNSECURE_SSL"`

	// Override region configuration
	ScalingoAPIURL  string `envconfig:"SCALINGO_API_URL"`
	ScalingoAuthURL string `envconfig:"SCALINGO_AUTH_URL"`
	ScalingoDbURL   string `envconfig:"SCALINGO_DB_URL"`
	ScalingoRegion  string `envconfig:"SCALINGO_REGION"`
	ScalingoSSHHost string `envconfig:"SCALINGO_SSH_HOST"`

	// Configuration files
	ConfigDir      string `envconfig:"CONFIG_DIR"`
	AuthFile       string `envconfig:"AUTH_FILE"`
	LogFile        string `envconfig:"LOG_FILE"`
	ConfigFilePath string `envconfig:"CONFIG_FILE_PATH"`
	ConfigFile     ConfigFile

	// Cache related files
	CacheDir         string `envconfig:"CACHE_DIR"`
	RegionsCachePath string `envconfig:"REGIONS_CACHE_PATH"`

	// Logging
	logFile *os.File
	Logger  *log.Logger
}

var (
	env = map[string]string{
		"SCALINGO_AUTH_URL":  "https://auth.scalingo.com",
		"SCALINGO_API_URL":   "",
		"SCALINGO_DB_URL":    "",
		"SCALINGO_SSH_HOST":  "",
		"SCALINGO_REGION":    "",
		"API_VERSION":        "1",
		"UNSECURE_SSL":       "false",
		"CONFIG_DIR":         ".config/scalingo",
		"CACHE_DIR":          ".cache/scalingo",
		"AUTH_FILE":          "auth",
		"CONFIG_FILE_PATH":   "config.json",
		"REGIONS_CACHE_PATH": "regions.json",
		"LOG_FILE":           "local.log",
	}
	C         Config
	TLSConfig *tls.Config
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

	rollbar.Token = RollbarToken
	rollbar.Platform = "client"
	rollbar.Environment = "production"
	rollbar.ErrorWriter = C.logFile

	TLSConfig = &tls.Config{}
	if C.UnsecureSsl {
		TLSConfig.InsecureSkipVerify = true
		TLSConfig.MinVersion = tls.VersionTLS10
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

func (config Config) CurrentUser(ctx context.Context) (*scalingo.User, error) {
	authenticator := &CliAuthenticator{}
	user, _, err := authenticator.LoadAuth(ctx)
	return user, err
}

type ClientConfigOpts struct {
	Region   string
	APIToken string
	AuthOnly bool
}

func (config Config) scalingoClientBaseConfig(opts ClientConfigOpts) scalingo.ClientConfig {
	return scalingo.ClientConfig{
		TLSConfig:    TLSConfig,
		Region:       opts.Region,
		AuthEndpoint: config.ScalingoAuthURL,
		APIToken:     opts.APIToken,
		UserAgent:    "Scalingo CLI v" + Version,
	}
}

func (config Config) scalingoClientConfig(ctx context.Context, opts ClientConfigOpts) (scalingo.ClientConfig, error) {
	c := config.scalingoClientBaseConfig(opts)
	if !opts.AuthOnly {
		if config.ScalingoAPIURL != "" && config.ScalingoDbURL != "" {
			c.APIEndpoint = config.ScalingoAPIURL
			c.DatabaseAPIEndpoint = config.ScalingoDbURL
		} else {
			region, err := GetRegion(ctx, config, config.ScalingoRegion, GetRegionOpts{
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

func ScalingoClientFromToken(ctx context.Context, token string) (*scalingo.Client, error) {
	config, err := C.scalingoClientConfig(ctx, ClientConfigOpts{APIToken: token})
	if err != nil {
		return nil, errgo.Notef(err, "fail to create Scalingo client")
	}
	return scalingo.New(ctx, config)
}

func ScalingoAuthClientFromToken(ctx context.Context, token string) (*scalingo.Client, error) {
	config := C.scalingoClientBaseConfig(ClientConfigOpts{
		AuthOnly: true, APIToken: token,
	})
	return scalingo.New(ctx, config)
}

func ScalingoAuthClient(ctx context.Context) (*scalingo.Client, error) {
	auth := &CliAuthenticator{}
	_, token, err := auth.LoadAuth(ctx)
	if err != nil {
		return nil, errgo.Notef(err, "fail to load authentication")
	}
	return ScalingoAuthClientFromToken(ctx, token.Token)
}

func ScalingoClient(ctx context.Context) (*scalingo.Client, error) {
	authenticator := &CliAuthenticator{}
	_, token, err := authenticator.LoadAuth(ctx)
	if err != nil {
		return nil, errgo.Notef(err, "fail to load credentials")
	}
	config, err := C.scalingoClientConfig(ctx, ClientConfigOpts{APIToken: token.Token})
	if err != nil {
		return nil, errgo.Notef(err, "fail to create Scalingo client")
	}
	return scalingo.New(ctx, config)
}

func ScalingoDatabaseClient(ctx context.Context) (*scalingo.PreviewClient, error) {
	authenticator := &CliAuthenticator{}
	_, token, err := authenticator.LoadAuth(ctx)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "load credentials")
	}
	config, err := C.scalingoClientConfig(ctx, ClientConfigOpts{APIToken: token.Token})
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create Scalingo client")
	}
	c, err := scalingo.New(ctx, config)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create Scalingo client")
	}
	return scalingo.NewPreviewClient(c), nil
}

func ScalingoClientForRegion(ctx context.Context, region string) (*scalingo.Client, error) {
	authenticator := &CliAuthenticator{}
	_, token, err := authenticator.LoadAuth(ctx)
	if err != nil {
		return nil, errgo.Notef(err, "fail to load credentials")
	}

	config := C.scalingoClientBaseConfig(ClientConfigOpts{
		Region:   region,
		APIToken: token.Token,
	})

	return scalingo.New(ctx, config)
}

func ScalingoUnauthenticatedAuthClient(ctx context.Context) (*scalingo.Client, error) {
	config := C.scalingoClientBaseConfig(ClientConfigOpts{AuthOnly: true})
	return scalingo.New(ctx, config)
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
