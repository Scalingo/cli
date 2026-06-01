package scalingo

import (
	"context"
	"net/http"
	"strconv"
	"time"

	httpclient "github.com/Scalingo/go-scalingo/v11/http"
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
var ErrOTPRequired = httpclient.ErrOTPRequired

// IsOTPRequired tests if the authentication backend return an OTP Required error
//
// Deprecated: use httpclient.IsOTPRequired instead of this wrapper.
func IsOTPRequired(err error) bool {
	return httpclient.IsOTPRequired(err)
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
	var btRes BearerTokenRes
	req := &httpclient.APIRequest{
		NoAuth:   true,
		Method:   http.MethodPost,
		Endpoint: "/tokens/exchange",
		Password: token,
	}

	err := c.AuthAPI().DoRequest(ctx, req, &btRes)
	if err != nil {
		return "", errors.Wrap(ctx, err, "make request POST /v1/tokens/exchange")
	}

	return btRes.Token, nil
}

func (c *Client) TokenCreateWithLogin(ctx context.Context, params TokenCreateParams, login LoginParams) (Token, error) {
	req := &httpclient.APIRequest{
		NoAuth:   true,
		Method:   http.MethodPost,
		Endpoint: "/tokens",
		Expected: httpclient.Statuses{http.StatusCreated},
		Username: login.Identifier,
		Password: login.Password,
		OTP:      login.OTP,
		Token:    login.JWT,
		Params:   map[string]any{"token": params},
	}

	var tokenRes TokenRes
	err := c.AuthAPI().DoRequest(ctx, req, &tokenRes)
	if err != nil {
		if httpclient.IsOTPRequired(err) {
			return Token{}, httpclient.ErrOTPRequired
		}
		return Token{}, errors.Wrap(ctx, err, "create token with login")
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
