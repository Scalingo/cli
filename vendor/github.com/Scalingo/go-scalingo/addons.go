package scalingo

import (
	"github.com/Scalingo/go-scalingo/http"

	"gopkg.in/errgo.v1"
)

type AddonsService interface {
	AddonsList(app string) ([]*Addon, error)
	AddonProvision(app, addon, planID string) (AddonRes, error)
	AddonDestroy(app, addonID string) error
	AddonUpgrade(app, addonID, planID string) (AddonRes, error)
	AddonToken(app, addonID string) (string, error)
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

type AddonToken struct {
	Token string `json:"token"`
}
type AddonTokenRes struct {
	Addon AddonToken `json:"addon"`
}

func (c *Client) AddonsList(app string) ([]*Addon, error) {
	var addonsRes AddonsRes
	err := c.ScalingoAPI().SubresourceList("apps", app, "addons", nil, &addonsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return addonsRes.Addons, nil
}

func (c *Client) AddonShow(app, addonID string) (Addon, error) {
	var addonRes AddonRes

	err := c.ScalingoAPI().SubresourceGet("apps", app, "addons", addonID, nil, &addonRes)
	if err != nil {
		return Addon{}, errgo.Mask(err, errgo.Any)
	}

	return addonRes.Addon, nil
}

func (c *Client) AddonProvision(app, addon, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "addons", AddonRes{Addon: Addon{AddonProviderID: addon, PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func (c *Client) AddonDestroy(app, addonID string) error {
	return c.ScalingoAPI().SubresourceDelete("apps", app, "addons", addonID)
}

func (c *Client) AddonUpgrade(app, addonID, planID string) (AddonRes, error) {
	var addonRes AddonRes
	err := c.ScalingoAPI().SubresourceUpdate("apps", app, "addons", addonID, AddonRes{Addon: Addon{PlanID: planID}}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func (c *Client) AddonToken(app, addonID string) (string, error) {
	var res AddonTokenRes
	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/addons/" + addonID + "/token",
	}, &res)
	if err != nil {
		return "", errgo.Notef(err, "fail to get addon token")
	}

	return res.Addon.Token, nil
}
