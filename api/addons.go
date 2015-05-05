package api

import "github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

type Addon struct {
	ID              string         `json:"id"`
	ResourceID      string         `json:"resource_id"`
	PlanID          string         `json:"plan_id"`
	AddonProviderID string         `json:"addon_provider_id"`
	Plan            *Plan          `json:"plan"`
	AddonProvider   *AddonProvider `json:"addon_provider"`
}

type AddonsRes struct {
	Addons []*Addon `json:"addons"`
}

type AddonRes struct {
	Addon     Addon    `json:"addon"`
	Message   string   `json:"message,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

func AddonsList(app string) ([]*Addon, error) {
	var addonsRes AddonsRes
	err := subresourceList(app, "addons", nil, &addonsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return addonsRes.Addons, nil
}

func AddonProvision(app, addon, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := subresourceAdd(app, "addons", AddonRes{Addon: Addon{AddonProviderID: addon, PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func AddonDestroy(app, addonID string) error {
	return subresourceDelete(app, "addons", addonID)
}

func AddonUpgrade(app, addonID, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := subresourceUpdate(app, "addons", addonID, AddonRes{Addon: Addon{PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}
