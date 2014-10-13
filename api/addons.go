package api

import (
	"encoding/json"
	"errors"
	"net/http"
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

func AddonsList() (*http.Response, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "GET",
		"endpoint": "/addons",
	}
	return Do(req)
}

func AddonPlansList(addon string) ([]*Plan, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "GET",
		"endpoint": "/addons/" + addon + "/plans",
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("addon " + addon + " not found.")
	}

	var params PlansParams
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return nil, err
	}

	return params.Plans, nil
}
