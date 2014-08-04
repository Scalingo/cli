package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Appsdeck/appsdeck/auth"
	"github.com/Appsdeck/appsdeck/config"
	"github.com/Appsdeck/appsdeck/debug"
	"github.com/Appsdeck/appsdeck/httpclient"
	"github.com/Appsdeck/appsdeck/session"
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

	// Manage authentication when not specified or when req is true
	if a, ok := req["auth"]; !ok || a.(bool) {
		if auth.Config.AuthToken != "" {
			params["user_email"] = auth.Config.Email
			params["user_token"] = auth.Config.AuthToken
		}
	}

	var urlWithoutParams string
	if u, ok := req["url"]; ok {
		urlWithoutParams = u.(string)
	} else {
		urlWithoutParams = fmt.Sprintf("%s%s", host, req["endpoint"].(string))
	}

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
	case "GET", "DELETE":
		values := url.Values{}
		for key, value := range params {
			values.Add(key, fmt.Sprintf("%v", value))
		}

		// In the case an url is given, we add the extra parameters
		var urlWithParams string
		if _, ok := req["url"]; ok {
			urlWithParams = fmt.Sprintf("%s&%s", urlWithoutParams, values.Encode())
		} else {
			urlWithParams = fmt.Sprintf("%s?%s", urlWithoutParams, values.Encode())
		}
		httpReq, err = http.NewRequest(req["method"].(string), urlWithParams, nil)
		if err != nil {
			return nil, err
		}
	}

	debug.Printf("[API] %v %v\n", httpReq.Method, httpReq.URL)
	if len(params) != 0 {
		debug.Printf("      Params : %v", params)
	}

	res, err := httpclient.Do(httpReq)
	if err != nil {
		fmt.Printf("Fail to query %s: %v\n", httpReq.Host, err)
		os.Exit(1)
	}

	if res.StatusCode == 500 {
		fmt.Fprintln(os.Stderr, "Server Internal Error - Our team has been contacted")
		os.Exit(1)
	}

	if res.StatusCode == 401 {
		fmt.Fprintln(os.Stderr, "You are not authorized to do this operation")
		session.DestroyToken()
		os.Exit(1)
	}

	return res, err
}
