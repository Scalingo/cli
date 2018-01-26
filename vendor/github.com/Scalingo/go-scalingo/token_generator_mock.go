package scalingo

import (
	"sync"
)

type tokenGeneratorMock struct {
	*sync.Mutex
	Calls         []interface{}
	accessToken   string
	responseError error
}

func newtokenGeneratorMock() *tokenGeneratorMock {
	return &tokenGeneratorMock{
		Mutex: &sync.Mutex{},
	}
}

func (c *tokenGeneratorMock) GetAccessToken() (string, error) {
	c.Lock()
	defer c.Unlock()

	c.Calls = append(c.Calls, nil)
	return c.accessToken, c.responseError
}

func (c *tokenGeneratorMock) SetClient(client *Client) {
	panic("not implemented")
}

func (c *tokenGeneratorMock) setAccessToken(t string) {
	c.Lock()
	defer c.Unlock()
	c.accessToken = t
}

func (c *tokenGeneratorMock) setResponseError(err error) {
	c.Lock()
	defer c.Unlock()
	c.responseError = err
}
