package http

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Scalingo/go-utils/errors/v3"
)

type TokensService interface {
	TokenExchange(ctx context.Context, token string) (string, error)
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

func (t *APITokenGenerator) GetAccessToken(ctx context.Context) (string, error) {
	// Ask for a new JWT if there wasn't any or if the current token will expire in less than 5 minutes
	if t.currentJWTexp.IsZero() || time.Until(t.currentJWTexp) < 5*time.Minute {
		jwtToken, err := t.TokensService.TokenExchange(ctx, t.APIToken)
		if err != nil {
			return "", errors.Wrap(ctx, err, "get access token")
		}

		token, _, err := jwt.NewParser().ParseUnverified(jwtToken, &apiJWTClaims{})
		if err != nil {
			return "", errors.Wrap(ctx, err, "parse JWT token")
		}

		if token.Valid {
			return "", errors.New(ctx, "JWT token is not valid")
		}

		if claims, ok := token.Claims.(*apiJWTClaims); ok {
			t.currentJWTexp = claims.ExpiresAt.Time
		} else {
			return "", errors.Errorf(ctx, "invalid exp date for jwt token: %v", token.Claims)
		}

		t.currentJWT = jwtToken
	}
	return t.currentJWT, nil
}
