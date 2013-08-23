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
)

func Do(req map[string]interface{}) (*http.Response, error) {
	host := config.C["APPSDECK_API"]
	if _, ok := req["host"]; ok {
		host = req["host"].(string)
	}

	params := make(map[string]interface{})
	if _, ok := req["params"]; ok {
		params = req["params"].(map[string]interface{})
	}

	if AuthToken != "" {
		params["auth_token"] = AuthToken
	}

	urlWithoutParams := fmt.Sprintf("http://%s%s", host, req["endpoint"].(string))

	var httpReq *http.Request
	var err error
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

	if AuthToken != "" {
		AddAuthToken(httpReq)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	debug.Printf("[API] %v %v\n", httpReq.Method, httpReq.URL)
	if len(params) != 0 {
		debug.Printf("      Params : %v", params)
	}

	return http.DefaultClient.Do(httpReq)
}

func AddAuthToken(req *http.Request) {
	req.Header.Set("Authorization", "Token token=\""+AuthToken+"\"")
}
