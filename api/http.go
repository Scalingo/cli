package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/cli/users"
)

var CurrentUser *users.User
var Prefix = config.C["API_PREFIX"]

type Statuses []int

func (statuses Statuses) Contains(status int) bool {
	for _, s := range statuses {
		if s == status {
			return true
		}
	}
	return false
}

// Execute an API request and return its response/error
func Do(req map[string]interface{}) (*http.Response, error) {
	var httpReq *http.Request
	var err error

	host := config.C["API_URL"]

	// If there is a host param, we use this custom host to send the request
	if _, ok := req["host"]; ok {
		host = req["host"].(string)
	}

	var params interface{}
	if _, ok := req["params"]; ok {
		params = req["params"]
	} else {
		params = make(map[string]interface{})
	}

	var urlWithoutParams string
	if u, ok := req["url"]; ok {
		urlWithoutParams = u.(string)
	} else {
		urlWithoutParams = fmt.Sprintf("%s%s%s", host, Prefix, req["endpoint"].(string))
	}

	// Execute the HTTP request according to the HTTP method
	switch req["method"].(string) {
	case "PATCH":
		fallthrough
	case "POST":
		fallthrough
	case "WITH_BODY":
		buffer, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(buffer)
		httpReq, err = http.NewRequest(req["method"].(string), urlWithoutParams, reader)
		if err != nil {
			return nil, err
		}
	case "GET", "DELETE":
		values := url.Values{}
		for key, value := range params.(map[string]interface{}) {
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

	// Manage authentication when not specified or when req is true
	a, ok := req["auth"]
	auth := !ok || a.(bool)
	if auth {
		user, err := AuthFromConfig()
		if err != nil {
			return nil, err
		}
		if user == nil {
			fmt.Println("You need to be authenticated to use Scalingo client.\nNo account ? → https://my.scalingo.com/users/signup")
			user, err = Auth()
			if err != nil {
				return nil, err
			}
		}
		CurrentUser = user
		httpReq.SetBasicAuth("", CurrentUser.AuthToken)
	}

	if token, ok := req["token"]; ok {
		httpReq.SetBasicAuth("", token.(string))
	}

	debug.Printf("[API] %v %v\n", httpReq.Method, httpReq.URL)
	debug.Printf(io.Indent(fmt.Sprintf("Headers: %v", httpReq.Header), 6))
	debug.Printf(io.Indent("Params : %v", 6), params)

	res, err := httpclient.Do(httpReq)
	if err != nil {
		fmt.Printf("Fail to query %s: %v\n", httpReq.Host, err)
		os.Exit(1)
	}

	if res.StatusCode == 500 {
		fmt.Fprintln(os.Stderr, "Server Internal Error - Our team has been contacted")
		os.Exit(1)
	}

	if auth && res.StatusCode == 401 {
		fmt.Fprintln(os.Stderr, "You are not authorized to do this operation")
		session.DestroyToken()
		os.Exit(1)
	}

	if _, ok := req["expected"]; ok && !req["expected"].(Statuses).Contains(res.StatusCode) {
		return nil, fmt.Errorf("Invalid status from server: %v", res.Status)
	}

	return res, err
}

func ParseJSON(res *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	debug.Println(string(body))

	return json.Unmarshal(body, data)
}
