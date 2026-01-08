package scalingo

import (
	"context"
	"strconv"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v9/http"
)

type AddonProvidersService interface {
	AddonProvidersList(context.Context) ([]*AddonProvider, error)
	AddonProviderPlansList(ctx context.Context, addon string, opts AddonProviderPlansListOpts) ([]*Plan, error)
}

var _ AddonProvidersService = (*Client)(nil)

type Plan struct {
	ID                        string `json:"id"`
	DisplayName               string `json:"display_name"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	Position                  int    `json:"position"`
	OnDemand                  bool   `json:"on_demand"`
	Disabled                  bool   `json:"disabled"`
	DisabledAlternativePlanID string `json:"disabled_alternative_plan_id"`
	SKU                       string `json:"sku"`
	HDSAvailable              bool   `json:"hds_available"`
	ToBeDiscontinued          bool   `json:"to_be_discontinued"`
	TrialAvailable            bool   `json:"trial_available"`
}

type AddonProviderPlansListResponse struct {
	Plans []*Plan `json:"plans"`
}

type AddonProviderPlansListOpts struct {
	ShowAll bool
}

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}

type AddonProvider struct {
	ID               string   `json:"id"`
	LogoURL          string   `json:"logo_url"`
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	Category         Category `json:"category"`
	ProviderName     string   `json:"provider_name"`
	ProviderURL      string   `json:"provider_url"`
	HDSAvailable     bool     `json:"hds_available"`
	Plans            []Plan   `json:"plans"`
}

type AddonProvidersListResponse struct {
	AddonProviders []*AddonProvider `json:"addon_providers"`
}

func (c *Client) AddonProvidersList(ctx context.Context) ([]*AddonProvider, error) {
	req := &http.APIRequest{
		NoAuth:   !c.isAuthenticatedClient(),
		Endpoint: "/addon_providers",
	}
	var response AddonProvidersListResponse
	err := c.ScalingoAPI().DoRequest(ctx, req, &response)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return response.AddonProviders, nil
}

var addonProviderTypo = map[string]string{
	"scalingo-mongo":    "scalingo-mongodb",
	"scalingo-influx":   "scalingo-influxdb",
	"scalingo-postgres": "scalingo-postgresql",
	"scalingo-postgre":  "scalingo-postgresql",
	"scalingo-pgsql":    "scalingo-postgresql",
	"scalingo-psql":     "scalingo-postgresql",
}

func (c *Client) AddonProviderPlansList(ctx context.Context, addon string, opts AddonProviderPlansListOpts) ([]*Plan, error) {
	correctAddon, ok := addonProviderTypo[addon]
	if ok {
		addon = correctAddon
	}

	params := map[string]string{
		"show_all": strconv.FormatBool(opts.ShowAll),
	}

	var response AddonProviderPlansListResponse
	req := &http.APIRequest{
		NoAuth:   !c.isAuthenticatedClient(),
		Endpoint: "/addon_providers/" + addon + "/plans",
		Params:   params,
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &response)
	if err != nil {
		return nil, errgo.Notef(err, "fail to get plans")
	}
	return response.Plans, nil
}
