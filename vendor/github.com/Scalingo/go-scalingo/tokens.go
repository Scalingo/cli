package scalingo

import (
	"fmt"
	"time"

	errgo "gopkg.in/errgo.v1"
)

type TokensService interface {
	TokensList() (Tokens, error)
	CreateToken(t Token) (Token, error)
	ShowToken(id int) (Token, error)
}

type TokensClient struct {
	*backendConfiguration
}

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

func (c *TokensClient) TokensList() (Tokens, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
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

func (c *TokensClient) CreateToken(t Token) (Token, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
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

func (c *TokensClient) ShowToken(id int) (Token, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		Endpoint: fmt.Sprintf("/tokens/%d", id),
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
