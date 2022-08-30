package scalingo

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/http"
)

type AddonProvidersService interface {
	AddonProvidersList(context.Context) ([]*AddonProvider, error)
	AddonProviderPlansList(ctx context.Context, addon string) ([]*Plan, error)
}

var _ AddonProvidersService = (*Client)(nil)

type Plan struct {
	ID          string  `json:"id"`
	LogoURL     string  `json:"logo_url"`
	DisplayName string  `json:"display_name"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}

type PlansParams struct {
	Plans []*Plan `json:"plans"`
}

type AddonProvider struct {
	ID      string `json:"id"`
	LogoURL string `json:"logo_url"`
	Name    string `json:"name"`
}

type ListParams struct {
	AddonProviders []*AddonProvider `json:"addon_providers"`
}

func (c *Client) AddonProvidersList(ctx context.Context) ([]*AddonProvider, error) {
	req := &http.APIRequest{
		NoAuth:   true,
		Endpoint: "/addon_providers",
	}
	var params ListParams
	err := c.ScalingoAPI().DoRequest(ctx, req, &params)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return params.AddonProviders, nil
}

var addonProviderTypo = map[string]string{
	"scalingo-mongo":    "scalingo-mongodb",
	"scalingo-influx":   "scalingo-influxdb",
	"scalingo-postgres": "scalingo-postgresql",
	"scalingo-postgre":  "scalingo-postgresql",
	"scalingo-pgsql":    "scalingo-postgresql",
	"scalingo-psql":     "scalingo-postgresql",
}

func (c *Client) AddonProviderPlansList(ctx context.Context, addon string) ([]*Plan, error) {
	correctAddon, ok := addonProviderTypo[addon]
	if ok {
		addon = correctAddon
	}

	var params PlansParams
	req := &http.APIRequest{
		NoAuth:   true,
		Endpoint: "/addon_providers/" + addon + "/plans",
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &params)
	if err != nil {
		return nil, errgo.Notef(err, "fail to get plans")
	}
	return params.Plans, nil
}
