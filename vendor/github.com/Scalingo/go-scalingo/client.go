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
	TokenGenerator TokenGenerator
	Endpoint       string
	TLSConfig      *tls.Config
	APIVersion     string
	httpClient     HTTPClient
}

type ClientConfig struct {
	Timeout        time.Duration
	Endpoint       string
	TLSConfig      *tls.Config
	TokenGenerator TokenGenerator
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
		TokenGenerator: cfg.TokenGenerator,
		Endpoint:       cfg.Endpoint,
		APIVersion:     defaultAPIVersion,
		TLSConfig:      cfg.TLSConfig,
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
