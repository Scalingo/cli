package http

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"
)

type AddonTokenGenerator struct {
	appID     string
	addonID   string
	exchanger AdddonTokenExchanger
}

type AdddonTokenExchanger interface {
	AddonToken(ctx context.Context, app, addonID string) (string, error)
}

func NewAddonTokenGenerator(app, addon string, exchanger AdddonTokenExchanger) TokenGenerator {
	return &AddonTokenGenerator{
		appID:     app,
		addonID:   addon,
		exchanger: exchanger,
	}
}

func (c *AddonTokenGenerator) GetAccessToken(ctx context.Context) (string, error) {
	token, err := c.exchanger.AddonToken(ctx, c.appID, c.addonID)
	if err != nil {
		return "", errors.Wrap(ctx, err, "get addon token")
	}
	return token, nil
}
