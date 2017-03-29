package scalingo

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	APIToken   string
	Endpoint   string
	TLSConfig  *tls.Config
	APIVersion string
	httpClient HTTPClient
}

type ClientConfig struct {
	Timeout   time.Duration
	Endpoint  string
	APIToken  string
	TLSConfig *tls.Config
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
	return &Client{
		APIToken:   cfg.APIToken,
		Endpoint:   cfg.Endpoint,
		APIVersion: defaultAPIVersion,
		TLSConfig:  cfg.TLSConfig,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Proxy:           http.ProxyFromEnvironment,
				TLSClientConfig: cfg.TLSConfig,
			},
		},
	}
}

func (c *Client) HTTPClient() HTTPClient {
	return c.httpClient
}
