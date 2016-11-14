package scalingo

import (
	"crypto/tls"
	"net/http"
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
	Endpoint  string
	APIToken  string
	TLSConfig *tls.Config
}

func NewClient(cfg ClientConfig) *Client {
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
