package api

import (
	"appsdeck/cli/config"
	"appsdeck/cli/debug"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	AuthToken string
	AuthEmail string
)

// Execute an API request and return its response/error
func Do(req map[string]interface{}) (*http.Response, error) {
	var httpReq *http.Request
	var err error

	host := config.C["APPSDECK_API"]

	// If there is a host param, we use this custom host to send the request
	if _, ok := req["host"]; ok {
		host = req["host"].(string)
	}

	params := make(map[string]interface{})
	if _, ok := req["params"]; ok {
		params = req["params"].(map[string]interface{})
	}

	// Manage authentication
	if AuthToken != "" {
		params["user_email"] = AuthEmail
		params["user_token"] = AuthToken
	}

	urlWithoutParams := fmt.Sprintf("%s%s", host, req["endpoint"].(string))

	// Execute the HTTP request according to the HTTP method
	switch req["method"].(string) {
	case "POST":
		buffer, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(buffer)
		httpReq, err = http.NewRequest("POST", urlWithoutParams, reader)
		if err != nil {
			return nil, err
		}
	case "GET":
		values := url.Values{}
		for key, value := range params {
			values.Add(key, value.(string))
		}
		urlWithParams := fmt.Sprintf("%s?%s", urlWithoutParams, values.Encode())
		httpReq, err = http.NewRequest("GET", urlWithParams, nil)
		if err != nil {
			return nil, err
		}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	debug.Printf("[API] %v %v\n", httpReq.Method, httpReq.URL)
	if len(params) != 0 {
		debug.Printf("      Params : %v", params)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config.TlsConfig,
		},
	}

	return client.Do(httpReq)
}
