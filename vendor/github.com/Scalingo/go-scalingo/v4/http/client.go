package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/errgo.v1"
)

const (
	AuthAPI     = "AUTHENTICATION_API"
	ScalingoAPI = "SCALINGO_API"
	DBAPI       = "DATABASES_API"
)

type APIConfig struct {
	Prefix string
}

type Client interface {
	ResourceList(resource string, payload, data interface{}) error
	ResourceAdd(resource string, payload, data interface{}) error
	ResourceGet(resource, resourceID string, payload, data interface{}) error
	ResourceUpdate(resource, resourceID string, payload, data interface{}) error
	ResourceDelete(resource, resourceID string) error

	SubresourceList(resource, resourceID, subresource string, payload, data interface{}) error
	SubresourceAdd(resource, resourceID, subresource string, payload, data interface{}) error
	SubresourceGet(resource, resourceID, subresource, id string, payload, data interface{}) error
	SubresourceUpdate(resource, resourceID, subresource, id string, payload, data interface{}) error
	SubresourceDelete(resource, resourceID, subresource, id string) error
	DoRequest(req *APIRequest, data interface{}) error
	Do(req *APIRequest) (*http.Response, error)

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
}

type client struct {
	tokenGenerator TokenGenerator
	endpoint       string
	userAgent      string
	apiConfig      APIConfig
	httpClient     *http.Client
	prefix         string
}

func NewClient(api string, cfg ClientConfig) Client {
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
	}

	return &c
}

func (c *client) ResourceGet(resource, resourceID string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource + "/" + resourceID,
		Params:   payload,
	}, data)
}

func (c *client) ResourceList(resource string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource,
		Params:   payload,
	}, data)
}

func (c *client) ResourceAdd(resource string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "POST",
		Endpoint: "/" + resource,
		Expected: Statuses{201},
		Params:   payload,
	}, data)
}

func (c client) ResourceUpdate(resource, resourceID string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "PATCH",
		Endpoint: "/" + resource + "/" + resourceID,
		Params:   payload,
	}, data)
}

func (c *client) ResourceDelete(resource, resourceID string) error {
	return c.DoRequest(&APIRequest{
		Method:   "DELETE",
		Endpoint: "/" + resource + "/" + resourceID,
		Expected: Statuses{204},
	}, nil)
}

func (c *client) SubresourceGet(resource, resourceID, subresource, id string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *client) SubresourceList(resource, resourceID, subresource string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "GET",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource,
		Params:   payload,
	}, data)
}

func (c *client) SubresourceAdd(resource, resourceID, subresource string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "POST",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource,
		Expected: Statuses{201},
		Params:   payload,
	}, data)
}

func (c *client) SubresourceDelete(resource, resourceID, subresource, id string) error {
	return c.DoRequest(&APIRequest{
		Method:   "DELETE",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource + "/" + id,
		Expected: Statuses{204},
	}, nil)
}

func (c *client) SubresourceUpdate(resource, resourceID, subresource, id string, payload, data interface{}) error {
	return c.DoRequest(&APIRequest{
		Method:   "PATCH",
		Endpoint: "/" + resource + "/" + resourceID + "/" + subresource + "/" + id,
		Params:   payload,
	}, data)
}

func (c *client) DoRequest(req *APIRequest, data interface{}) error {
	if c.endpoint == "" {
		return errgo.New("API Endpoint is not defined, did you forget to pass the Region field to the New method?")
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if data == nil {
		return nil
	}

	err = parseJSON(res, data)
	if err != nil {
		return errgo.Notef(err, "fail to parse JSON of subresource response")
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
