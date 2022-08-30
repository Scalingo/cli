package http

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/errgo.v1"
)

type TokensService interface {
	TokenExchange(token string) (string, error)
}

type APITokenGenerator struct {
	APIToken      string
	TokensService TokensService

	// Internal cache of JWT info
	currentJWT    string
	currentJWTexp time.Time
}

type apiJWTClaims struct {
	jwt.RegisteredClaims
}

func NewAPITokenGenerator(tokensService TokensService, apiToken string) *APITokenGenerator {
	return &APITokenGenerator{
		APIToken:      apiToken,
		TokensService: tokensService,
	}
}

func (t *APITokenGenerator) GetAccessToken() (string, error) {
	// Ask for a new JWT if there wasn't any or if the current token will expire in less than 5 minutes
	if t.currentJWTexp.IsZero() || t.currentJWTexp.Sub(time.Now()) < 5*time.Minute {
		jwtToken, err := t.TokensService.TokenExchange(t.APIToken)
		if err != nil {
			return "", errgo.Notef(err, "fail to get access token")
		}
		token, err := jwt.ParseWithClaims(jwtToken, &apiJWTClaims{}, nil)
		// If token is nil, nothing has been parsed, if it's not, err will be a
		// ValidatingError we want to ignore
		if token == nil {
			return "", errgo.Notef(err, "fail to parse jwt token")
		}

		if claims, ok := token.Claims.(*apiJWTClaims); ok {
			t.currentJWTexp = claims.ExpiresAt.Time
		} else {
			return "", errgo.Notef(err, "invalid exp date for jwt token: %v", token.Claims)
		}

		t.currentJWT = jwtToken
	}
	return t.currentJWT, nil
}
