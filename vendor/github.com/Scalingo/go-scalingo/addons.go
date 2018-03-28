package scalingo

import "gopkg.in/errgo.v1"

type AddonsService interface {
	AddonsList(app string) ([]*Addon, error)
	AddonProvision(app, addon, planID string) (AddonRes, error)
	AddonDestroy(app, addonID string) error
	AddonUpgrade(app, addonID, planID string) (AddonRes, error)
}

type AddonsClient struct {
	subresourceClient
}

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

func (c *AddonsClient) AddonsList(app string) ([]*Addon, error) {
	var addonsRes AddonsRes
	err := c.subresourceList(app, "addons", nil, &addonsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return addonsRes.Addons, nil
}

func (c *AddonsClient) AddonProvision(app, addon, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := c.subresourceAdd(app, "addons", AddonRes{Addon: Addon{AddonProviderID: addon, PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func (c *AddonsClient) AddonDestroy(app, addonID string) error {
	return c.subresourceDelete(app, "addons", addonID)
}

func (c *AddonsClient) AddonUpgrade(app, addonID, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := c.subresourceUpdate(app, "addons", addonID, AddonRes{Addon: Addon{PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}
