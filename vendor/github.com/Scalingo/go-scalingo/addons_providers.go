package scalingo

import (
	"encoding/json"

	"gopkg.in/errgo.v1"
)

type AddonProvidersService interface {
	AddonProvidersList() ([]*AddonProvider, error)
	AddonProviderPlansList(addon string) ([]*Plan, error)
}

type AddonProvidersClient struct {
	*backendConfiguration
}

type Plan struct {
	ID               string `json:"id"`
	LogoURL          string `json:"logo_url"`
	DisplayName      string `json:"display_name"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
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

func (c *AddonProvidersClient) AddonProvidersList() ([]*AddonProvider, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		NoAuth:   true,
		Endpoint: "/addon_providers",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var params ListParams
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
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

func (c *AddonProvidersClient) AddonProviderPlansList(addon string) ([]*Plan, error) {
	correctAddon, ok := addonProviderTypo[addon]
	if ok {
		addon = correctAddon
	}

	req := &APIRequest{
		Client:   c.backendConfiguration,
		NoAuth:   true,
		Endpoint: "/addon_providers/" + addon + "/plans",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params PlansParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return params.Plans, nil
}
