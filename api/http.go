package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/httpclient"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/users"
)

var CurrentUser *users.User

type APIRequest struct {
	NoAuth      bool
	URL         string
	Method      string
	Endpoint    string
	Token       string
	Expected    Statuses
	Params      interface{}
	HTTPRequest *http.Request
}

type Statuses []int

func (req *APIRequest) FillDefaultValues() error {
	if req.URL == "" {
		req.URL = fmt.Sprintf("%s%s", ApiUrl, config.C.ApiPrefix)
	}
	if req.Method == "" {
		req.Method = "GET"
	}
	if req.Expected == nil || len(req.Expected) == 0 {
		req.Expected = Statuses{200}
	}
	if req.Params == nil {
		req.Params = make(map[string]interface{})
	}
	if req.Token == "" && !req.NoAuth {
		user, err := AuthFromConfig()
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		if user == nil {
			fmt.Println("You need to be authenticated to use Scalingo client.\nNo account ? → https://scalingo.com")
			user, err = Auth()
			if err != nil {
				return errgo.Mask(err, errgo.Any)
			}
		}
		CurrentUser = user
		req.Token = CurrentUser.AuthToken
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
	debug.Printf(io.Indent("Params : %v", 6), req.Params)

	req.HTTPRequest.SetBasicAuth("", req.Token)
	res, err := httpclient.Do(req.HTTPRequest)
	if err != nil {
		fmt.Printf("Fail to query %s: %v\n", req.HTTPRequest.Host, err)
		os.Exit(1)
	}

	if req.Expected.Contains(res.StatusCode) {
		return res, nil
	}

	return nil, NewRequestFailedError(res, req)
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
