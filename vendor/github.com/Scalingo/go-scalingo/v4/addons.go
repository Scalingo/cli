package scalingo

import (
	"encoding/json"
	"time"

	"github.com/Scalingo/go-scalingo/v4/http"

	"gopkg.in/errgo.v1"
)

type AddonsService interface {
	AddonsList(app string) ([]*Addon, error)
	AddonProvision(app string, params AddonProvisionParams) (AddonRes, error)
	AddonDestroy(app, addonID string) error
	AddonUpgrade(app, addonID string, params AddonUpgradeParams) (AddonRes, error)
	AddonToken(app, addonID string) (string, error)
	AddonLogsURL(app, addonID string) (string, error)
}

var _ AddonsService = (*Client)(nil)

type AddonStatus string

const (
	AddonStatusRunning      AddonStatus = "running"
	AddonStatusProvisioning AddonStatus = "provisioning"
	AddonStatusSuspended    AddonStatus = "suspended"
)

type Addon struct {
	ID              string         `json:"id"`
	AppID           string         `json:"app_id"`
	ResourceID      string         `json:"resource_id"`
	Status          AddonStatus    `json:"status"`
	Plan            *Plan          `json:"plan"`
	AddonProvider   *AddonProvider `json:"addon_provider"`
	ProvisionedAt   time.Time      `json:"provisioned_at"`
	DeprovisionedAt time.Time      `json:"deprovisioned_at"`
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

type AddonLogsURLRes struct {
	URL string `json:"url"`
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

// AddonProvisionParams gathers all arguments which can be sent to provision an addon
type AddonProvisionParams struct {
	AddonProviderID string            `json:"addon_provider_id"`
	PlanID          string            `json:"plan_id"`
	Options         map[string]string `json:"options"`
}

type AddonProvisionParamsWrapper struct {
	Addon AddonProvisionParams `json:"addon"`
}

func (c *Client) AddonProvision(app string, params AddonProvisionParams) (AddonRes, error) {
	var addonRes AddonRes
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "addons", AddonProvisionParamsWrapper{params}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func (c *Client) AddonDestroy(app, addonID string) error {
	return c.ScalingoAPI().SubresourceDelete("apps", app, "addons", addonID)
}

type AddonUpgradeParams struct {
	PlanID string `json:"plan_id"`
}

type AddonUpgradeParamsWrapper struct {
	Addon AddonUpgradeParams `json:"addon"`
}

func (c *Client) AddonUpgrade(app, addonID string, params AddonUpgradeParams) (AddonRes, error) {
	var addonRes AddonRes
	err := c.ScalingoAPI().SubresourceUpdate(
		"apps", app, "addons", addonID,
		AddonUpgradeParamsWrapper{Addon: params}, &addonRes,
	)
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

func (c *Client) AddonLogsURL(app, addonID string) (string, error) {
	var url AddonLogsURLRes
	res, err := c.DBAPI(app, addonID).Do(&http.APIRequest{
		Endpoint: "/databases/" + addonID + "/logs",
	})
	if err != nil {
		return "", errgo.Notef(err, "fail to get log URL")
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&url)
	if err != nil {
		return "", errgo.Notef(err, "invalid response")
	}

	return url.URL, nil
}
