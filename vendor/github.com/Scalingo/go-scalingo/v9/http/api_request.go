package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v9/debug"
	pkgio "github.com/Scalingo/go-scalingo/v9/io"
)

type APIRequest struct {
	NoAuth      bool
	URL         string
	Method      string
	Endpoint    string
	Expected    Statuses
	Params      interface{}
	HTTPRequest *http.Request
	Token       string // Directly use a Bearer token
	Username    string // Username for the OAuth generator (nil if you use a token)
	Password    string // Password for the OAuth generator
	OTP         string // OTP value
}

type Statuses []int

func (c *client) fillDefaultValues(ctx context.Context, req *APIRequest) error {
	if req.Method == "" {
		req.Method = "GET"
	}
	if req.Expected == nil || len(req.Expected) == 0 {
		req.Expected = Statuses{200}
	}
	if req.Params == nil {
		req.Params = make(map[string]interface{})
	}

	if !req.NoAuth && c.IsAuthenticatedClient() {
		var err error
		req.Token, err = c.TokenGenerator().GetAccessToken(ctx)
		if err != nil {
			return errgo.Notef(err, "fail to get the access token for this request")
		}
	}

	if req.URL == "" {
		req.URL = c.BaseURL()
	}
	return nil
}

func (statuses Statuses) Contains(status int) bool {
	for _, s := range statuses {
		if s == status {
			return true
		}
	}
	return false
}

// Execute an API request and return its response/error
func (c *client) Do(ctx context.Context, req *APIRequest) (*http.Response, error) {
	err := c.fillDefaultValues(ctx, req)
	if err != nil {
		return nil, errgo.Notef(err, "fail to fill client with default values")
	}

	endpoint := req.URL + req.Endpoint

	// Execute the HTTP request according to the HTTP method
	var body io.Reader
	switch req.Method {
	case http.MethodPatch, http.MethodPost, http.MethodPut, "WITH_BODY":
		buffer, err := json.Marshal(req.Params)
		if err != nil {
			return nil, errgo.Notef(err, "fail to marshal params")
		}
		body = bytes.NewReader(buffer)
	case http.MethodGet, http.MethodDelete:
		values, err := req.BuildQueryFromParams()
		if err != nil {
			return nil, errgo.Notef(err, "fail to build the query params")
		}
		endpoint = fmt.Sprintf("%s?%s", endpoint, values.Encode())
	}

	req.HTTPRequest, err = http.NewRequest(req.Method, endpoint, body)
	if err != nil {
		return nil, errgo.Notef(err, "fail to initialize the '%s' query", req.Method)
	}
	req.HTTPRequest.Header.Add("User-Agent", c.userAgent)
	requestID, ok := ctx.Value("request_id").(string)
	if ok {
		req.HTTPRequest.Header.Add("X-Request-ID", requestID)
	}

	for key, values := range c.extraHeaders {
		for _, value := range values {
			req.HTTPRequest.Header.Add(key, value)
		}
	}

	debug.Printf("[API] %v %v\n", req.HTTPRequest.Method, req.HTTPRequest.URL)
	debug.Printf("%s", pkgio.Indent(fmt.Sprintf("User Agent: %v", req.HTTPRequest.UserAgent()), 6))
	debug.Printf("%s", pkgio.Indent(fmt.Sprintf("Headers: %v", req.HTTPRequest.Header), 6))
	debug.Printf(pkgio.Indent("Params: %v", 6), req.Params)

	if req.Token != "" {
		req.HTTPRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", req.Token))
	} else if req.Username != "" || req.Password != "" {
		req.HTTPRequest.SetBasicAuth(req.Username, req.Password)
	}

	if req.OTP != "" {
		req.HTTPRequest.Header.Add("X-Authorization-OTP", req.OTP)
	}

	now := time.Now()
	res, err := c.doRequest(req.HTTPRequest)
	if err != nil {
		return nil, fmt.Errorf("Fail to query %s: %v", req.HTTPRequest.Host, err)
	}
	debug.Printf(pkgio.Indent("Request ID: %v", 6), res.Header.Get("X-Request-Id"))
	debug.Printf(pkgio.Indent("Duration: %v", 6), time.Since(now))

	if req.Expected.Contains(res.StatusCode) {
		return res, nil
	}

	return nil, NewRequestFailedError(res, req)
}

func (c *client) doRequest(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("User-Agent", "Scalingo Go Client")
	return c.HTTPClient().Do(req)
}

func parseJSON(res *http.Response, data interface{}) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		reqErr := fmt.Sprintf("%v %v", res.Request.Method, res.Request.URL)
		return errgo.Newf("fail to read body of request %v: %v", reqErr, err)
	}

	debug.Println(string(body))

	err = json.Unmarshal(body, data)
	if err != nil {
		reqErr := fmt.Sprintf("%v %v", res.Request.Method, res.Request.URL)
		return errgo.Newf("fail to parse JSON of request %v: %v", reqErr, err)
	}

	return nil
}

func (req *APIRequest) BuildQueryFromParams() (url.Values, error) {
	values := url.Values{}
	if reflect.TypeOf(req.Params).Kind() != reflect.Map {
		return nil, errgo.Newf("%#v is not a map", req.Params)
	}
	if reflect.TypeOf(req.Params).Key().Kind() != reflect.String {
		return nil, errgo.Newf("%#v is not a map of string", req.Params)
	}
	value := reflect.ValueOf(req.Params)
	for _, key := range value.MapKeys() {
		value := value.MapIndex(key)
		values.Add(fmt.Sprintf("%v", key), fmt.Sprintf("%v", value.Interface()))
	}
	return values, nil
}
