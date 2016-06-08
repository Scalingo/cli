package scalingo

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	APIToken   string
	Endpoint   string
	TLSConfig  *tls.Config
	APIVersion string
	httpClient *http.Client
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

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}
