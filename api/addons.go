package api

import "github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

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
	req := &APIRequest{
		Endpoint: "/apps/" + app + "/addons",
	}
	res, err := req.Do()
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
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/addons",
		Expected: Statuses{201},
		Params: map[string]interface{}{
			"addon": map[string]interface{}{
				"addon_provider_id": addon,
				"plan_id":           planID,
			},
		},
	}
	res, err := req.Do()
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
	req := &APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/addons/" + addonID,
		Expected: Statuses{204},
	}
	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	return nil
}

func AddonUpgrade(app, addonID, planID string) (*UpgradeAddonParams, error) {
	req := &APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/addons" + addonID,
		Params: map[string]interface{}{
			"addon": map[string]interface{}{
				"plan_id": planID,
			},
		},
	}
	res, err := req.Do()
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
