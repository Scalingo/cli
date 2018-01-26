package test

import (
	"net/http"
	"sync"
)

type HTTPClient struct {
	*sync.Mutex
	Calls         []*http.Request
	responseData  *http.Response
	responseError error
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Mutex: &sync.Mutex{},
	}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.Lock()
	defer c.Unlock()

	c.Calls = append(c.Calls, req)
	return c.responseData, c.responseError
}

func (c *HTTPClient) SetResponseData(res *http.Response) {
	c.Lock()
	defer c.Unlock()
	c.responseData = res
}

func (c *HTTPClient) SetResponseError(err error) {
	c.Lock()
	defer c.Unlock()
	c.responseError = err
}
