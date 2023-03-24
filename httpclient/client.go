package httpclient

import (
	"net/http"

	"github.com/Scalingo/cli/config"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: config.TLSConfig,
		},
	}
)

func Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("User-Agent", "Scalingo CLI v"+config.Version)
	return client.Do(req)
}
