package scalingo

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v9/http"
)

type API interface {
	AddonsService
	AddonProvidersService
	AppsService
	AlertsService
	AutoscalersService
	BackupsService
	CollaboratorsService
	ContainersService
	ContainerSizesService
	CronTasksService
	DatabasesService
	DeploymentsService
	DomainsService
	EventsService
	InvoicesService
	KeysService
	LogDrainsService
	LogsArchivesService
	LogsService
	NotificationPlatformsService
	NotifiersService
	OperationsService
	PrivateNetworksService
	RegionsService
	RunsService
	SignUpService
	SourcesService
	StacksService
	TokensService
	UsersService
	VariablesService
	httpclient.TokenGenerator

	ScalingoAPI() httpclient.Client
	AuthAPI() httpclient.Client
	DBAPI(app, addon string) httpclient.Client
}

var _ API = (*Client)(nil)

type Client struct {
	config     ClientConfig
	apiClient  httpclient.Client
	dbClient   httpclient.Client
	authClient httpclient.Client
}

type ClientConfig struct {
	Timeout                time.Duration
	TLSConfig              *tls.Config
	APIEndpoint            string
	APIPrefix              string
	AuthEndpoint           string
	AuthPrefix             string
	DatabaseAPIEndpoint    string
	DatabaseAPIPrefix      string
	APIToken               string
	Region                 string
	UserAgent              string
	DisableHTTPClientCache bool
	ExtraHeaders           ExtraHeaders

	// StaticTokenGenerator is present for Scalingo internal use only
	StaticTokenGenerator *StaticTokenGenerator
}

type ExtraHeaders struct {
	API         http.Header
	DatabaseAPI http.Header
	Auth        http.Header
}

func New(ctx context.Context, cfg ClientConfig) (*Client, error) {
	// Apply defaults
	if cfg.AuthEndpoint == "" {
		cfg.AuthEndpoint = "https://auth.scalingo.com"
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "go-scalingo v" + Version
	}

	// If there's no region defined return the client as is
	if cfg.Region == "" {
		return &Client{
			config: cfg,
		}, nil
	}

	// if a region was defined, create a temp client to query the auth service for region list
	// then create the real client
	tmpClient := &Client{
		config: cfg,
	}

	region, err := tmpClient.getRegion(ctx, cfg.Region)
	if err == ErrRegionNotFound {
		return nil, err
	} else if err != nil {
		return nil, errgo.Notef(err, "fail to get region informations")
	}

	cfg.APIEndpoint = region.API
	cfg.DatabaseAPIEndpoint = region.DatabaseAPI
	return &Client{
		config: cfg,
	}, nil
}

func (c *Client) ScalingoAPI() httpclient.Client {
	if c.apiClient != nil {
		return c.apiClient
	}

	var tokenGenerator httpclient.TokenGenerator
	if c.config.StaticTokenGenerator != nil {
		tokenGenerator = c.config.StaticTokenGenerator
	}
	if c.config.APIToken != "" {
		tokenGenerator = httpclient.NewAPITokenGenerator(c, c.config.APIToken)
	}
	prefix := "/v1"
	if c.config.APIPrefix != "" {
		prefix = c.config.APIPrefix
	}

	client := httpclient.NewClient(httpclient.ClientConfig{
		UserAgent:      c.config.UserAgent,
		Timeout:        c.config.Timeout,
		TLSConfig:      c.config.TLSConfig,
		APIConfig:      httpclient.APIConfig{Prefix: prefix},
		TokenGenerator: tokenGenerator,
		Endpoint:       c.config.APIEndpoint,
		ExtraHeaders:   c.config.ExtraHeaders.API,
	})

	if !c.config.DisableHTTPClientCache {
		c.apiClient = client
	}

	return client
}

func (c *Client) DBAPI(app, addon string) httpclient.Client {
	if c.dbClient != nil {
		return c.dbClient
	}

	prefix := "/api"
	if c.config.DatabaseAPIPrefix != "" {
		prefix = c.config.DatabaseAPIPrefix
	}
	return httpclient.NewClient(httpclient.ClientConfig{
		UserAgent:      c.config.UserAgent,
		Timeout:        c.config.Timeout,
		TLSConfig:      c.config.TLSConfig,
		APIConfig:      httpclient.APIConfig{Prefix: prefix},
		TokenGenerator: httpclient.NewAddonTokenGenerator(app, addon, c),
		Endpoint:       c.config.DatabaseAPIEndpoint,
		ExtraHeaders:   c.config.ExtraHeaders.DatabaseAPI,
	})
}

func (c *Client) AuthAPI() httpclient.Client {
	if c.authClient != nil {
		return c.authClient
	}

	var tokenGenerator httpclient.TokenGenerator
	if c.config.StaticTokenGenerator != nil {
		tokenGenerator = c.config.StaticTokenGenerator
	}
	if c.config.APIToken != "" {
		tokenGenerator = httpclient.NewAPITokenGenerator(c, c.config.APIToken)
	}

	prefix := "/v1"
	if c.config.AuthPrefix != "" {
		prefix = c.config.AuthPrefix
	}
	client := httpclient.NewClient(httpclient.ClientConfig{
		UserAgent:      c.config.UserAgent,
		Timeout:        c.config.Timeout,
		TLSConfig:      c.config.TLSConfig,
		APIConfig:      httpclient.APIConfig{Prefix: prefix},
		TokenGenerator: tokenGenerator,
		Endpoint:       c.config.AuthEndpoint,
		ExtraHeaders:   c.config.ExtraHeaders.Auth,
	})

	if !c.config.DisableHTTPClientCache {
		c.authClient = client
	}

	return client
}

// Preview returns the databases next generation APIs client, available as preview feature.
func (c *Client) Preview() *PreviewClient {
	return NewPreviewClient(c)
}

func (c *Client) isAuthenticatedClient() bool {
	return c.ScalingoAPI().IsAuthenticatedClient()
}
