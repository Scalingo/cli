package scalingo

import (
	"fmt"
	"time"

	errgo "gopkg.in/errgo.v1"
)

type TokensService interface {
	TokensList() (Tokens, error)
	TokenCreate(TokenCreateParams) (Token, error)
	TokenExchange(TokenExchangeParams) (string, error)
	TokenShow(id int) (Token, error)
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

type TokenCreateParams struct {
	Name string `json:"name"`
}

type TokenExchangeParams struct {
	Token string
}

type Tokens []*Token

type TokensRes struct {
	Tokens Tokens `json:"tokens"`
}

type BearerTokenRes struct {
	Token string `json:"token"`
}

type TokenRes struct {
	Token Token `json:"token"`
}

func (c *TokensClient) TokensList() (Tokens, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		URL:      AuthURL(),
		Endpoint: "/v1/tokens",
	}

	res, err := req.Do()
	if err != nil {
		return nil, errgo.Notef(err, "fail to get tokens")
	}

	var tokensRes TokensRes
	err = ParseJSON(res, &tokensRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to parse response from server")
	}

	return tokensRes.Tokens, nil
}

func (c *TokensClient) TokenExchange(params TokenExchangeParams) (string, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		NoAuth:   true,
		Method:   "POST",
		URL:      AuthURL(),
		Endpoint: "/v1/tokens/exchange",
		Password: params.Token,
	}

	res, err := req.Do()
	if err != nil {
		return "", errgo.Notef(err, "fail to make request POST /v1/tokens/exchange")
	}

	var btRes BearerTokenRes
	err = ParseJSON(res, &btRes)
	if err != nil {
		return "", errgo.NoteMask(err, "invalid response from authentication service", errgo.Any)
	}

	return btRes.Token, nil
}

func (c *TokensClient) TokenCreateWithLogin(params TokenCreateParams, login LoginParams) (Token, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		NoAuth:   true,
		Method:   "POST",
		URL:      AuthURL(),
		Endpoint: "/v1/tokens",
		Expected: Statuses{201},
		Username: login.Identifier,
		Password: login.Password,
		OTP:      login.OTP,
		Token:    login.JWT,
		Params:   map[string]interface{}{"token": params},
	}

	resp, err := req.Do()
	if err != nil {
		if IsOTPRequired(err) {
			return Token{}, ErrOTPRequired
		}
		return Token{}, errgo.Notef(err, "request failed")
	}

	var tokenRes TokenRes
	err = ParseJSON(resp, &tokenRes)
	if err != nil {
		return Token{}, errgo.NoteMask(err, "invalid response from authentication service", errgo.Any)
	}

	return tokenRes.Token, nil
}

func (c *TokensClient) TokenCreate(params TokenCreateParams) (Token, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		URL:      AuthURL(),
		Expected: Statuses{201},
		Endpoint: "/v1/tokens",
		Params:   map[string]interface{}{"token": params},
		Method:   "POST",
	}

	res, err := req.Do()
	if err != nil {
		return Token{}, errgo.Notef(err, "fail to create token")
	}

	var tokenRes TokenRes
	err = ParseJSON(res, &tokenRes)
	if err != nil {
		return Token{}, errgo.Notef(err, "fail to parse response from server")
	}

	return tokenRes.Token, nil
}

func (c *TokensClient) TokenShow(id int) (Token, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		URL:      AuthURL(),
		Endpoint: fmt.Sprintf("/v1/tokens/%d", id),
	}

	res, err := req.Do()
	if err != nil {
		return Token{}, errgo.Notef(err, "fail to get token")
	}

	var tokenRes TokenRes
	err = ParseJSON(res, &tokenRes)
	if err != nil {
		return Token{}, errgo.Notef(err, "fail to parse response from server")
	}

	return tokenRes.Token, nil
}
