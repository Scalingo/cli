package api

import "net/http"

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
		"expected": Statuses{200},
	}
	return Do(req)
}

func AddonPlansList(addon string) ([]*Plan, error) {
	req := map[string]interface{}{
		"auth":     false,
		"method":   "GET",
		"endpoint": "/addons/" + addon + "/plans",
		"expected": Statuses{200},
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var params PlansParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, err
	}

	return params.Plans, nil
}
