package scalingo

import (
	"errors"
	"os"

	"gopkg.in/errgo.v1"
)

type TokenGenerator interface {
	GetAccessToken() (string, error)
}

type LoginParams struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	OTP        string `json:"otp"`
	JWT        string `json:"jwt"`
}

type APITokenGenerator struct {
	APIToken      string        `json:"token"`
	TokensService TokensService `json:"-"`
}

var ErrOTPRequired = errors.New("OTP Required")

// IsOTPRequired tests if the authentication backend return an OTP Required error
func IsOTPRequired(err error) bool {
	rerr, ok := err.(*RequestFailedError)
	if !ok {
		return false
	}

	if rerr.Message == "OTP Required" {
		return true
	}
	return false
}

const defaultAuthUrl = "https://auth.scalingo.com"

func (c *Client) GetAPITokenGenerator(apiToken string) *APITokenGenerator {
	return &APITokenGenerator{
		APIToken:      apiToken,
		TokensService: c,
	}
}

func (t *APITokenGenerator) GetAccessToken() (string, error) {
	accessToken, err := t.TokensService.TokenExchange(TokenExchangeParams{Token: t.APIToken})
	if err != nil {
		return "", errgo.Notef(err, "fail to get access token")
	}
	return accessToken, nil
}

func (c *Client) AuthURL() string {
	if len(c.authEndpoint) != 0 {
		return c.authEndpoint
	}

	if os.Getenv("SCALINGO_AUTH_URL") != "" {
		c.authEndpoint = os.Getenv("SCALINGO_AUTH_URL")
	} else {
		c.authEndpoint = defaultAuthUrl
	}
	return c.authEndpoint
}
