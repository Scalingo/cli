package http

import (
	errgo "gopkg.in/errgo.v1"
)

type TokensService interface {
	TokenExchange(token string) (string, error)
}

type APITokenGenerator struct {
	APIToken      string        `json:"token"`
	TokensService TokensService `json:"-"`
}

func NewAPITokenGenerator(tokensService TokensService, apiToken string) TokenGenerator {
	return &APITokenGenerator{
		APIToken:      apiToken,
		TokensService: tokensService,
	}
}

func (t *APITokenGenerator) GetAccessToken() (string, error) {
	accessToken, err := t.TokensService.TokenExchange(t.APIToken)
	if err != nil {
		return "", errgo.Notef(err, "fail to get access token")
	}
	return accessToken, nil
}
