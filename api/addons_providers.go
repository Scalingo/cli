package api

import (
	"encoding/json"

	"gopkg.in/errgo.v1"
)

type Plan struct {
	ID          string `json:"id"`
	LogoURL     string `json:"logo_url"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PlansParams struct {
	Plans []*Plan `json:"plans"`
}

type AddonProvider struct {
	LogoURL   string `json:"logo_url"`
	Name      string `json:"name"`
	NameParam string `json:"name_param"`
}

type ListParams struct {
	AddonProviders []*AddonProvider `json:"addon_providers"`
}

func AddonProvidersList() ([]*AddonProvider, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "GET",
		"endpoint": "/addon-providers",
		"expected": Statuses{200},
	}
	res, err := Do(req)
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

func AddonProviderPlansList(addon string) ([]*Plan, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "GET",
		"endpoint": "/addon-providers/" + addon + "/plans",
		"expected": Statuses{200},
	}
	res, err := Do(req)
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
