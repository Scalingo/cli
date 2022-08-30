package scalingo

import (
	"context"
	"crypto/tls"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/http"
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
	DeploymentsService
	DomainsService
	VariablesService
	EventsService
	KeysService
	LogDrainsService
	LogsArchivesService
	LogsService
	NotificationPlatformsService
	NotifiersService
	OperationsService
	RegionsService
	RegionMigrationsService
	RunsService
	SignUpService
	SourcesService
	StacksService
	TokensService
	UsersService
	http.TokenGenerator

	ScalingoAPI() http.Client
	AuthAPI() http.Client
	DBAPI(app, addon string) http.Client
}

var _ API = (*Client)(nil)

type Client struct {
	config     ClientConfig
	apiClient  http.Client
	dbClient   http.Client
	authClient http.Client
}

type ClientConfig struct {
	Timeout             time.Duration
	TLSConfig           *tls.Config
	APIEndpoint         string
	APIPrefix           string
	AuthEndpoint        string
	AuthPrefix          string
	DatabaseAPIEndpoint string
	DatabaseAPIPrefix   string
	APIToken            string
	Region              string
	UserAgent           string

	// StaticTokenGenerator is present for Scalingo internal use only
	StaticTokenGenerator *StaticTokenGenerator
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

func (c *Client) ScalingoAPI() http.Client {
	if c.apiClient != nil {
		return c.apiClient
	}

	var tokenGenerator http.TokenGenerator
	if c.config.StaticTokenGenerator != nil {
		tokenGenerator = c.config.StaticTokenGenerator
	}
	if len(c.config.APIToken) != 0 {
		tokenGenerator = http.NewAPITokenGenerator(c, c.config.APIToken)
	}
	prefix := "/v1"
	if c.config.APIPrefix != "" {
		prefix = c.config.APIPrefix
	}

	return http.NewClient(http.ScalingoAPI, http.ClientConfig{
		UserAgent:      c.config.UserAgent,
		Timeout:        c.config.Timeout,
		TLSConfig:      c.config.TLSConfig,
		APIConfig:      http.APIConfig{Prefix: prefix},
		TokenGenerator: tokenGenerator,
		Endpoint:       c.config.APIEndpoint,
	})
}

func (c *Client) DBAPI(app, addon string) http.Client {
	if c.dbClient != nil {
		return c.dbClient
	}

	prefix := "/api"
	if c.config.DatabaseAPIPrefix != "" {
		prefix = c.config.DatabaseAPIPrefix
	}
	return http.NewClient(http.DBAPI, http.ClientConfig{
		UserAgent:      c.config.UserAgent,
		Timeout:        c.config.Timeout,
		TLSConfig:      c.config.TLSConfig,
		APIConfig:      http.APIConfig{Prefix: prefix},
		TokenGenerator: http.NewAddonTokenGenerator(app, addon, c),
		Endpoint:       c.config.DatabaseAPIEndpoint,
	})
}

func (c *Client) AuthAPI() http.Client {
	if c.authClient != nil {
		return c.authClient
	}

	var tokenGenerator http.TokenGenerator
	if c.config.StaticTokenGenerator != nil {
		tokenGenerator = c.config.StaticTokenGenerator
	}
	if len(c.config.APIToken) != 0 {
		tokenGenerator = http.NewAPITokenGenerator(c, c.config.APIToken)
	}

	prefix := "/v1"
	if c.config.AuthPrefix != "" {
		prefix = c.config.AuthPrefix
	}
	return http.NewClient(http.AuthAPI, http.ClientConfig{
		UserAgent:      c.config.UserAgent,
		Timeout:        c.config.Timeout,
		TLSConfig:      c.config.TLSConfig,
		APIConfig:      http.APIConfig{Prefix: prefix},
		TokenGenerator: tokenGenerator,
		Endpoint:       c.config.AuthEndpoint,
	})
}
