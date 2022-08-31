package scalingo

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/http"
)

type AddonsService interface {
	AddonsList(ctx context.Context, app string) ([]*Addon, error)
	AddonProvision(ctx context.Context, app string, params AddonProvisionParams) (AddonRes, error)
	AddonDestroy(ctx context.Context, app, addonID string) error
	AddonUpgrade(ctx context.Context, app, addonID string, params AddonUpgradeParams) (AddonRes, error)
	AddonToken(ctx context.Context, app, addonID string) (string, error)
	AddonLogsURL(ctx context.Context, app, addonID string) (string, error)
	AddonLogsArchives(ctx context.Context, app, addonID string, page int) (*LogsArchivesResponse, error)
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

func (c *Client) AddonsList(ctx context.Context, app string) ([]*Addon, error) {
	var addonsRes AddonsRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "addons", nil, &addonsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return addonsRes.Addons, nil
}

func (c *Client) AddonShow(ctx context.Context, app, addonID string) (Addon, error) {
	var addonRes AddonRes

	err := c.ScalingoAPI().SubresourceGet(ctx, "apps", app, "addons", addonID, nil, &addonRes)
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

func (c *Client) AddonProvision(ctx context.Context, app string, params AddonProvisionParams) (AddonRes, error) {
	var addonRes AddonRes
	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "addons", AddonProvisionParamsWrapper{params}, &addonRes)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func (c *Client) AddonDestroy(ctx context.Context, app, addonID string) error {
	return c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "addons", addonID)
}

type AddonUpgradeParams struct {
	PlanID string `json:"plan_id"`
}

type AddonUpgradeParamsWrapper struct {
	Addon AddonUpgradeParams `json:"addon"`
}

func (c *Client) AddonUpgrade(ctx context.Context, app, addonID string, params AddonUpgradeParams) (AddonRes, error) {
	var addonRes AddonRes
	err := c.ScalingoAPI().SubresourceUpdate(
		ctx, "apps", app, "addons", addonID,
		AddonUpgradeParamsWrapper{Addon: params}, &addonRes,
	)
	if err != nil {
		return AddonRes{}, errgo.Mask(err, errgo.Any)
	}
	return addonRes, nil
}

func (c *Client) AddonToken(ctx context.Context, app, addonID string) (string, error) {
	var res AddonTokenRes
	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/addons/" + addonID + "/token",
	}, &res)
	if err != nil {
		return "", errgo.Notef(err, "fail to get addon token")
	}

	return res.Addon.Token, nil
}

func (c *Client) AddonLogsURL(ctx context.Context, app, addonID string) (string, error) {
	var url AddonLogsURLRes
	res, err := c.DBAPI(app, addonID).Do(ctx, &http.APIRequest{
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

func (c *Client) AddonLogsArchives(ctx context.Context, app, addonID string, page int) (*LogsArchivesResponse, error) {
	res, err := c.DBAPI(app, addonID).Do(ctx, &http.APIRequest{
		Endpoint: "/databases/" + addonID + "/logs_archives",
		Params: map[string]string{
			"page": strconv.FormatInt(int64(page), 10),
		},
	})
	if err != nil {
		return nil, errgo.Notef(err, "fail to get log archives")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errgo.Notef(err, "fail to read body of response")
	}

	var logsRes = LogsArchivesResponse{}
	err = json.Unmarshal(body, &logsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to parse response")
	}

	return &logsRes, nil
}
