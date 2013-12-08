package httpclient

import (
	"appsdeck/config"
	"net/http"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config.TlsConfig,
		},
	}
)

func Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("User-Agent", "Appsdeck CLI v"+config.Version)
	return client.Do(req)
}
