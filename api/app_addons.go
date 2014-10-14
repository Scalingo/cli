package api

import "errors"

type AddonResource struct {
	ID         string `json:"id"`
	ResourceID string `json:"resource_id"`
	Plan       string `json:"plan"`
	PlanID     string `json:"plan_id"`
	Addon      string `json:"addon"`
	AddonID    string `json:"addon_id"`
}

type ListAddonResourcesParams struct {
	AddonResources []*AddonResource `json:"addon_resources"`
}

type ProvisionAddonResourceParams struct {
	AddonResource *AddonResource `json:"addon_resource"`
	Message       string         `json:"message"`
	Variables     []string       `json:"variables"`
}

type UpgradeAddonResourceParams ProvisionAddonResourceParams

func AddonResourcesList(app string) ([]*AddonResource, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/addons",
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var params ListAddonResourcesParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, err
	}
	return params.AddonResources, nil
}

func AddonResourceProvision(app, addon, planID string) (*ProvisionAddonResourceParams, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/addons",
		"params": map[string]interface{}{
			"addon_id": addon,
			"plan_id":  planID,
		},
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return nil, errors.New("invalid status code from server: " + res.Status)
	}

	var params *ProvisionAddonResourceParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func AddonResourceDestroy(app, addonResourceID string) error {
	req := map[string]interface{}{
		"method":   "DELETE",
		"endpoint": "/apps/" + app + "/addons/" + addonResourceID,
	}
	res, err := Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return errors.New("invalid status code fomr server: " + res.Status)
	}

	return nil
}

func AddonResourceUpgrade(app, addonResourceID, planID string) (*UpgradeAddonResourceParams, error) {
	req := map[string]interface{}{
		"method":   "PATCH",
		"endpoint": "/apps/" + app + "/addons/" + addonResourceID,
		"params": map[string]interface{}{
			"plan_id": planID,
		},
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("invalid status code from server: " + res.Status)
	}

	var params *UpgradeAddonResourceParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, err
	}

	return params, nil
}
