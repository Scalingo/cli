package scalingo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/Scalingo/go-scalingo/debug"
	"github.com/Scalingo/go-scalingo/io"
	"gopkg.in/errgo.v1"
)

var (
	defaultEndpoint   = "https://api.scalingo.com"
	defaultAPIVersion = "1"
	ErrNoAuth         = errgo.New("authentication required")
)

type APIRequest struct {
	Client      API
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

func (req *APIRequest) FillDefaultValues() error {
	if req.Method == "" {
		req.Method = "GET"
	}
	if req.Expected == nil || len(req.Expected) == 0 {
		req.Expected = Statuses{200}
	}
	if req.Params == nil {
		req.Params = make(map[string]interface{})
	}
	if req.Client == nil {
		req.Client = NewClient(ClientConfig{
			Endpoint:   defaultEndpoint,
			APIVersion: defaultAPIVersion,
		})
	}

	if !req.NoAuth {
		var err error
		req.Token, err = req.Client.GetAccessToken()
		if err != nil {
			return ErrNoAuth
		}
	}

	if req.URL == "" {
		req.URL = fmt.Sprintf("%s%s%s", req.Client.Endpoint(), "/v", req.Client.APIVersion())
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
func (req *APIRequest) Do() (*http.Response, error) {
	err := req.FillDefaultValues()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	endpoint := req.URL + req.Endpoint

	// Execute the HTTP request according to the HTTP method
	switch req.Method {
	case "PATCH":
		fallthrough
	case "POST":
		fallthrough
	case "PUT":
		fallthrough
	case "WITH_BODY":
		buffer, err := json.Marshal(req.Params)
		if err != nil {
			return nil, errgo.Mask(err, errgo.Any)
		}
		reader := bytes.NewReader(buffer)
		req.HTTPRequest, err = http.NewRequest(req.Method, endpoint, reader)
		if err != nil {
			return nil, errgo.Mask(err, errgo.Any)
		}
	case "GET", "DELETE":
		values, err := req.BuildQueryFromParams()
		if err != nil {
			return nil, errgo.Mask(err, errgo.Any)
		}
		endpoint = fmt.Sprintf("%s?%s", endpoint, values.Encode())
		req.HTTPRequest, err = http.NewRequest(req.Method, endpoint, nil)
		if err != nil {
			return nil, errgo.Mask(err, errgo.Any)
		}
	}

	debug.Printf("[API] %v %v\n", req.HTTPRequest.Method, req.HTTPRequest.URL)
	debug.Printf(io.Indent(fmt.Sprintf("Headers: %v", req.HTTPRequest.Header), 6))
	debug.Printf(io.Indent("Params: %v", 6), req.Params)

	if req.Token != "" {
		req.HTTPRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", req.Token))
	} else if req.Username != "" || req.Password != "" {
		req.HTTPRequest.SetBasicAuth(req.Username, req.Password)
	}

	if req.OTP != "" {
		req.HTTPRequest.Header.Add("X-Authorization-OTP", req.OTP)
	}

	now := time.Now()
	res, err := req.doRequest(req.HTTPRequest)
	if err != nil {
		fmt.Printf("Fail to query %s: %v\n", req.HTTPRequest.Host, err)
		os.Exit(1)
	}
	debug.Printf(io.Indent("Duration: %v", 6), time.Now().Sub(now))

	if req.Expected.Contains(res.StatusCode) {
		return res, nil
	}

	return nil, NewRequestFailedError(res, req)
}

func (apiReq *APIRequest) doRequest(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("User-Agent", "Scalingo Go Client")
	return apiReq.Client.HTTPClient().Do(req)
}

func ParseJSON(res *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errgo.Newf("fail to read body of request %v, %v", res.Request, err)
	}

	debug.Println(string(body))

	err = json.Unmarshal(body, data)
	if err != nil {
		return errgo.Newf("fail to parse JSON of request %v, %v", res.Request, err)
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
