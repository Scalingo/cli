package api

import "gopkg.in/errgo.v1"

type Addon struct {
	ID            string         `json:"id"`
	ResourceID    string         `json:"resource_id"`
	Plan          *Plan          `json:"plan"`
	AddonProvider *AddonProvider `json:"addon_provider"`
}

type ListAddonsParams struct {
	Addons []*Addon `json:"addons"`
}

type ProvisionAddonParams struct {
	Addon     *Addon   `json:"addon"`
	Message   string   `json:"message,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

type UpgradeAddonParams ProvisionAddonParams

func AddonsList(app string) ([]*Addon, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/addons",
		"expected": Statuses{200},
	}
	res, err := Do(req)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params ListAddonsParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return params.Addons, nil
}

func AddonProvision(app, addon, planID string) (*ProvisionAddonParams, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/addons",
		"params": map[string]interface{}{
			"addon": map[string]interface{}{
				"addon_provider_id": addon,
				"plan_id":           planID,
			},
		},
		"expected": Statuses{201},
	}
	res, err := Do(req)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params *ProvisionAddonParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return params, nil
}

func AddonDestroy(app, addonID string) error {
	req := map[string]interface{}{
		"method":   "DELETE",
		"endpoint": "/apps/" + app + "/addons/" + addonID,
		"expected": Statuses{204},
	}
	res, err := Do(req)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	return nil
}

func AddonUpgrade(app, addonID, planID string) (*UpgradeAddonParams, error) {
	req := map[string]interface{}{
		"method":   "PATCH",
		"endpoint": "/apps/" + app + "/addons/" + addonID,
		"params": map[string]interface{}{
			"addon": map[string]interface{}{
				"plan_id": planID,
			},
		},
		"expected": Statuses{200},
	}
	res, err := Do(req)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params *UpgradeAddonParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return params, nil
}
