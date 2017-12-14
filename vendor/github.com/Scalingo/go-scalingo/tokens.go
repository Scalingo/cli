package scalingo

import (
	"fmt"
	"time"

	errgo "gopkg.in/errgo.v1"
)

type Token struct {
	ID        int       `json:"int"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

type Tokens []*Token

type TokensResp struct {
	Tokens Tokens `json:"tokens"`
}

type TokenResp struct {
	Token Token `json:"token"`
}

func (c *Client) TokensList() (Tokens, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/tokens",
	}

	res, err := req.Do()

	if err != nil {
		return nil, errgo.Notef(err, "fail to get tokens")
	}

	var tokens TokensResp
	err = ParseJSON(res, &tokens)
	if err != nil {
		return nil, errgo.Notef(err, "fail to parse response from server")
	}

	return tokens.Tokens, nil
}

func (c *Client) CreateToken(t Token) (Token, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/tokens",
		Method:   "POST",
		Params:   t,
	}

	var token TokenResp

	res, err := req.Do()
	if err != nil {
		return token.Token, errgo.Notef(err, "fail to create token")
	}

	err = ParseJSON(res, &token)
	if err != nil {
		return token.Token, errgo.Notef(err, "fail to parse response from server")
	}

	return token.Token, nil
}

func (c *Client) ShowToken(id int) (Token, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: fmt.Sprintf("/tokens/%s", id),
	}

	var token TokenResp

	res, err := req.Do()
	if err != nil {
		return token.Token, errgo.Notef(err, "fail to get token")
	}

	err = ParseJSON(res, &token)
	if err != nil {
		return token.Token, errgo.Notef(err, "fail to parse response from server")
	}

	return token.Token, nil
}
