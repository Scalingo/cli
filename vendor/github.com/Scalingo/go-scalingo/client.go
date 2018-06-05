package scalingo

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type API interface {
	AddonsService
	AddonProvidersService
	AppsService
	AlertsService
	AutoscalersService
	CollaboratorsService
	DeploymentsService
	DomainsService
	VariablesService
	EventsService
	KeysService
	LoginService
	LogsArchivesService
	LogsService
	NotificationPlatformsService
	NotificationsService
	NotifiersService
	OperationsService
	RunsService
	SignUpService
	SourcesService
	TokensService
	UsersService

	TokenGenerator

	APIVersion() string
	Endpoint() string
	HTTPClient() HTTPClient
}

var _ API = (*Client)(nil)

type Client struct {
	tokenGenerator TokenGenerator
	endpoint       string
	authEndpoint   string
	TLSConfig      *tls.Config
	apiVersion     string
	httpClient     HTTPClient
}

type ClientConfig struct {
	Timeout        time.Duration
	Endpoint       string
	AuthEndpoint   string
	TLSConfig      *tls.Config
	TokenGenerator TokenGenerator
	APIVersion     string
	APIToken       string
}

func NewClient(cfg ClientConfig) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = defaultEndpoint
	}
	if cfg.TLSConfig == nil {
		cfg.TLSConfig = &tls.Config{}
	}

	if cfg.APIVersion == "" {
		cfg.APIVersion = defaultAPIVersion
	}

	c := Client{
		tokenGenerator: cfg.TokenGenerator,
		endpoint:       cfg.Endpoint,
		authEndpoint:   cfg.AuthEndpoint,
		apiVersion:     cfg.APIVersion,
		TLSConfig:      cfg.TLSConfig,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Proxy:           http.ProxyFromEnvironment,
				TLSClientConfig: cfg.TLSConfig,
			},
		},
	}

	if len(cfg.APIToken) != 0 && c.tokenGenerator == nil {
		c.tokenGenerator = c.GetAPITokenGenerator(cfg.APIToken)
	}

	return &c
}

func (c *Client) HTTPClient() HTTPClient {
	return c.httpClient
}

func (c *Client) GetAccessToken() (string, error) {
	if c.tokenGenerator == nil {
		return "", errors.New("no token generator")
	}
	return c.tokenGenerator.GetAccessToken()
}

func (c *Client) APIVersion() string {
	return c.apiVersion
}

func (c *Client) Endpoint() string {
	return c.endpoint
}
