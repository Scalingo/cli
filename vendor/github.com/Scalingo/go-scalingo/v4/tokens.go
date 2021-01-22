package scalingo

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Scalingo/go-scalingo/v4/http"
	errgo "gopkg.in/errgo.v1"
)

type TokensService interface {
	TokensList() (Tokens, error)
	TokenCreate(TokenCreateParams) (Token, error)
	TokenExchange(token string) (string, error)
	TokenShow(id int) (Token, error)
}

var _ TokensService = (*Client)(nil)
var ErrOTPRequired = errors.New("OTP Required")

// IsOTPRequired tests if the authentication backend return an OTP Required error
func IsOTPRequired(err error) bool {
	rerr, ok := err.(*http.RequestFailedError)
	if !ok {
		return false
	}

	if rerr.Message == "OTP Required" {
		return true
	}
	return false
}

type Token struct {
	ID        int       `json:"int"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

type LoginParams struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	OTP        string `json:"otp"`
	JWT        string `json:"jwt"`
}

type TokenCreateParams struct {
	Name string `json:"name"`
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

func (c *Client) TokensList() (Tokens, error) {
	var tokensRes TokensRes

	err := c.AuthAPI().ResourceList("tokens", nil, &tokensRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to get tokens")
	}

	return tokensRes.Tokens, nil
}

func (c *Client) TokenExchange(token string) (string, error) {
	req := &http.APIRequest{
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/tokens/exchange",
		Password: token,
	}

	res, err := c.AuthAPI().Do(req)
	if err != nil {
		return "", errgo.Notef(err, "fail to make request POST /v1/tokens/exchange")
	}
	defer res.Body.Close()

	var btRes BearerTokenRes
	err = json.NewDecoder(res.Body).Decode(&btRes)
	if err != nil {
		return "", errgo.Notef(err, "invalid response from authentication service")
	}

	return btRes.Token, nil
}

func (c *Client) TokenCreateWithLogin(params TokenCreateParams, login LoginParams) (Token, error) {
	req := &http.APIRequest{
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/tokens",
		Expected: http.Statuses{201},
		Username: login.Identifier,
		Password: login.Password,
		OTP:      login.OTP,
		Token:    login.JWT,
		Params:   map[string]interface{}{"token": params},
	}

	resp, err := c.AuthAPI().Do(req)
	if err != nil {
		if IsOTPRequired(err) {
			return Token{}, ErrOTPRequired
		}
		return Token{}, errgo.Notef(err, "request failed")
	}
	defer resp.Body.Close()

	var tokenRes TokenRes
	err = json.NewDecoder(resp.Body).Decode(&tokenRes)
	if err != nil {
		return Token{}, errgo.NoteMask(err, "invalid response from authentication service", errgo.Any)
	}

	return tokenRes.Token, nil
}

func (c *Client) TokenCreate(params TokenCreateParams) (Token, error) {
	var tokenRes TokenRes
	payload := map[string]TokenCreateParams{
		"token": params,
	}
	err := c.AuthAPI().ResourceAdd("tokens", payload, &tokenRes)
	if err != nil {
		return Token{}, errgo.Notef(err, "fail to create token")
	}

	return tokenRes.Token, nil
}

func (c *Client) TokenShow(id int) (Token, error) {
	var tokenRes TokenRes
	err := c.AuthAPI().ResourceGet("tokens", strconv.Itoa(id), nil, &tokenRes)
	if err != nil {
		return Token{}, errgo.Notef(err, "fail to get token")
	}

	return tokenRes.Token, nil
}

func (c *Client) GetAccessToken() (string, error) {
	return c.ScalingoAPI().TokenGenerator().GetAccessToken()
}
