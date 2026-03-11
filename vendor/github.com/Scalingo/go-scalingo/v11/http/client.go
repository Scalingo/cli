package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/Scalingo/go-utils/errors/v3"
)

type APIConfig struct {
	Prefix string
}

type Client interface {
	ResourceList(ctx context.Context, resource string, payload, data any) error
	ResourceAdd(ctx context.Context, resource string, payload, data any) error
	ResourceGet(ctx context.Context, resource, resourceID string, payload, data any) error
	ResourceUpdate(ctx context.Context, resource, resourceID string, payload, data any) error
	ResourceDelete(ctx context.Context, resource, resourceID string) error

	SubresourceList(ctx context.Context, resource, resourceID, subresource string, payload, data any) error
	SubresourceAdd(ctx context.Context, resource, resourceID, subresource string, payload, data any) error
	SubresourceGet(ctx context.Context, resource, resourceID, subresource, id string, payload, data any) error
	SubresourceUpdate(ctx context.Context, resource, resourceID, subresource, id string, payload, data any) error
	SubresourceDelete(ctx context.Context, resource, resourceID, subresource, id string) error
	DoRequest(ctx context.Context, req *APIRequest, data any) error
	Do(ctx context.Context, req *APIRequest) (*http.Response, error)

	TokenGenerator() TokenGenerator
	IsAuthenticatedClient() bool
	BaseURL() string
	HTTPClient() *http.Client
}

type ClientConfig struct {
	UserAgent      string
	Timeout        time.Duration
	TLSConfig      *tls.Config
	APIConfig      APIConfig
	Endpoint       string
	TokenGenerator TokenGenerator
	ExtraHeaders   http.Header
}

type client struct {
	tokenGenerator TokenGenerator
	endpoint       string
	userAgent      string
	apiConfig      APIConfig
	httpClient     *http.Client
	prefix         string
	extraHeaders   http.Header
}

func NewClient(cfg ClientConfig) Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.TLSConfig == nil {
		cfg.TLSConfig = &tls.Config{}
	}

	c := client{
		prefix:         cfg.APIConfig.Prefix,
		endpoint:       cfg.Endpoint,
		tokenGenerator: cfg.TokenGenerator,
		userAgent:      cfg.UserAgent,
		apiConfig:      cfg.APIConfig,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Proxy:           http.ProxyFromEnvironment,
				TLSClientConfig: cfg.TLSConfig,
			},
		},
		extraHeaders: cfg.ExtraHeaders,
	}

	return &c
}

func (c *client) ResourceGet(ctx context.Context, resource, resourceID string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource + "/" + resourceID,
		Params:   payload,
	}, data)
}

func (c *client) ResourceList(ctx context.Context, resource string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource,
		Params:   payload,
	}, data)
}

func (c *client) ResourceAdd(ctx context.Context, resource string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "POST",
		Endpoint: "/" + resource,
		Expected: Statuses{http.StatusCreated},
		Params:   payload,
	}, data)
}

func (c client) ResourceUpdate(ctx context.Context, resource, resourceID string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "PATCH",
		Endpoint: "/" + resource + "/" + resourceID,
		Params:   payload,
	}, data)
}

func (c *client) ResourceDelete(ctx context.Context, resource, resourceID string) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "DELETE",
		Endpoint: "/" + resource + "/" + resourceID,
		Expected: Statuses{http.StatusNoContent},
	}, nil)
}

func (c *client) SubresourceGet(ctx context.Context, resource, resourceID, subresource, id string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *client) SubresourceList(ctx context.Context, resource, resourceID, subresource string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource,
		Params:   payload,
	}, data)
}

func (c *client) SubresourceAdd(ctx context.Context, resource, resourceID, subresource string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "POST",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource,
		Expected: Statuses{http.StatusCreated},
		Params:   payload,
	}, data)
}

func (c *client) SubresourceDelete(ctx context.Context, resource, resourceID, subresource, id string) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "DELETE",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource + "/" + id,
		Expected: Statuses{http.StatusNoContent},
	}, nil)
}

func (c *client) SubresourceUpdate(ctx context.Context, resource, resourceID, subresource, id string, payload, data any) error {
	return c.DoRequest(ctx, &APIRequest{
		Method:   "PATCH",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *client) DoRequest(ctx context.Context, req *APIRequest, data any) error {
	if c.endpoint == "" {
		return errors.New(ctx, "API Endpoint is not defined, did you forget to pass the Region field to the New method?")
	}

	res, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if data == nil {
		return nil
	}

	_, err = parseJSON(ctx, res, data)
	if err != nil {
		return errors.Wrap(ctx, err, "parse JSON of subresource response")
	}
	return nil
}

func (c *client) IsAuthenticatedClient() bool {
	return c.tokenGenerator != nil
}

func (c *client) TokenGenerator() TokenGenerator {
	return c.tokenGenerator
}

func (c *client) BaseURL() string {
	endpoint := c.endpoint

	if c.prefix != "" {
		return fmt.Sprintf("%s%s", endpoint, c.prefix)
	}
	return endpoint
}

func (c *client) HTTPClient() *http.Client {
	return c.httpClient
}
