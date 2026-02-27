package scalingo

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/Scalingo/go-scalingo/v9/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type TokensService interface {
	TokensList(context.Context) (Tokens, error)
	TokenCreate(context.Context, TokenCreateParams) (Token, error)
	TokenExchange(ctx context.Context, token string) (string, error)
	TokenShow(ctx context.Context, id int) (Token, error)
}

var _ TokensService = (*Client)(nil)

// Deprecated: use http.ErrOTPRequired instead of this wrapper.
var ErrOTPRequired = http.ErrOTPRequired

// IsOTPRequired tests if the authentication backend return an OTP Required error
//
// Deprecated: use http.IsOTPRequired instead of this wrapper.
func IsOTPRequired(err error) bool {
	return http.IsOTPRequired(err)
}

type Token struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
	Token      string    `json:"token"`
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

func (c *Client) TokensList(ctx context.Context) (Tokens, error) {
	var tokensRes TokensRes

	err := c.AuthAPI().ResourceList(ctx, "tokens", nil, &tokensRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get tokens")
	}

	return tokensRes.Tokens, nil
}

func (c *Client) TokenExchange(ctx context.Context, token string) (string, error) {
	req := &http.APIRequest{
		NoAuth:   true,
		Method:   "POST",
		Endpoint: "/tokens/exchange",
		Password: token,
	}

	res, err := c.AuthAPI().Do(ctx, req)
	if err != nil {
		return "", errors.Wrap(ctx, err, "make request POST /v1/tokens/exchange")
	}
	defer res.Body.Close()

	var btRes BearerTokenRes
	err = json.NewDecoder(res.Body).Decode(&btRes)
	if err != nil {
		return "", errors.Wrap(ctx, err, "invalid response from authentication service")
	}

	return btRes.Token, nil
}

func (c *Client) TokenCreateWithLogin(ctx context.Context, params TokenCreateParams, login LoginParams) (Token, error) {
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

	resp, err := c.AuthAPI().Do(ctx, req)
	if err != nil {
		if http.IsOTPRequired(err) {
			return Token{}, http.ErrOTPRequired
		}
		return Token{}, errors.Wrap(ctx, err, "create token with login")
	}
	defer resp.Body.Close()

	var tokenRes TokenRes
	err = json.NewDecoder(resp.Body).Decode(&tokenRes)
	if err != nil {
		return Token{}, errors.Wrap(ctx, err, "invalid response from authentication service")
	}

	return tokenRes.Token, nil
}

func (c *Client) TokenCreate(ctx context.Context, params TokenCreateParams) (Token, error) {
	var tokenRes TokenRes
	payload := map[string]TokenCreateParams{
		"token": params,
	}
	err := c.AuthAPI().ResourceAdd(ctx, "tokens", payload, &tokenRes)
	if err != nil {
		return Token{}, errors.Wrap(ctx, err, "create token")
	}

	return tokenRes.Token, nil
}

func (c *Client) TokenShow(ctx context.Context, id int) (Token, error) {
	var tokenRes TokenRes
	err := c.AuthAPI().ResourceGet(ctx, "tokens", strconv.Itoa(id), nil, &tokenRes)
	if err != nil {
		return Token{}, errors.Wrap(ctx, err, "get token")
	}

	return tokenRes.Token, nil
}

func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	return c.ScalingoAPI().TokenGenerator().GetAccessToken(ctx)
}
