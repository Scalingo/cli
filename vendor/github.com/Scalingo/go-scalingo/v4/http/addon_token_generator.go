package http

import (
	"gopkg.in/errgo.v1"
)

type AddonTokenGenerator struct {
	appID     string
	addonID   string
	exchanger AdddonTokenExchanger
}

type AdddonTokenExchanger interface {
	AddonToken(app, addonID string) (string, error)
}

func NewAddonTokenGenerator(app, addon string, exchanger AdddonTokenExchanger) TokenGenerator {
	return &AddonTokenGenerator{
		appID:     app,
		addonID:   addon,
		exchanger: exchanger,
	}
}

func (c *AddonTokenGenerator) GetAccessToken() (string, error) {
	token, err := c.exchanger.AddonToken(c.appID, c.addonID)
	if err != nil {
		return "", errgo.Notef(err, "fail to get addon token")
	}
	return token, nil
}
