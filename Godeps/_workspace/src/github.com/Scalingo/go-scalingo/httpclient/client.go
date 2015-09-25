package httpclient

import (
	"crypto/tls"
	"net/http"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{},
		},
	}
)

func Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return client.Do(req)
}
