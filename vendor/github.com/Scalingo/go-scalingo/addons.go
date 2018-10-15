package scalingo

import "gopkg.in/errgo.v1"

type AddonsService interface {
	AddonsList(app string) ([]*Addon, error)
	AddonProvision(app, addon, planID string) (AddonRes, error)
	AddonDestroy(app, addonID string) error
	AddonUpgrade(app, addonID, planID string) (AddonRes, error)
}

var _ AddonsService = (*Client)(nil)

type Addon struct {
	ID              string         `json:"id"`
	AppID           string         `json:"app_id"`
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

func (c *Client) AddonsList(app string) ([]*Addon, error) {
	var addonsRes AddonsRes
	err := c.subresourceList(app, "addons", nil, &addonsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return addonsRes.Addons, nil
}

func (c *Client) AddonShow(app, addonID string) (Addon, error) {
	var addonRes AddonRes

	err := c.subresourceGet(app, "addons", addonID, nil, &addonRes)
	if err != nil {
		return Addon{}, errgo.Mask(err, errgo.Any)
	}

	return addonRes.Addon, nil
}

func (c *Client) AddonProvision(app, addon, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := c.subresourceAdd(app, "addons", AddonRes{Addon: Addon{AddonProviderID: addon, PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func (c *Client) AddonDestroy(app, addonID string) error {
	return c.subresourceDelete(app, "addons", addonID)
}

func (c *Client) AddonUpgrade(app, addonID, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := c.subresourceUpdate(app, "addons", addonID, AddonRes{Addon: Addon{PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}
