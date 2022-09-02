package scalingo

import (
	"context"
	"net/http"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v5/http"
)

type LogDrainsService interface {
	LogDrainsList(ctx context.Context, app string) ([]LogDrain, error)
	LogDrainAdd(ctx context.Context, app string, params LogDrainAddParams) (*LogDrainRes, error)
	LogDrainRemove(ctx context.Context, app, URL string) error
	LogDrainAddonRemove(ctx context.Context, app, addonID string, URL string) error
	LogDrainsAddonList(ctx context.Context, app string, addonID string) ([]LogDrain, error)
	LogDrainAddonAdd(ctx context.Context, app string, addonID string, params LogDrainAddParams) (*LogDrainRes, error)
}

var _ LogDrainsService = (*Client)(nil)

type LogDrain struct {
	AppID string `json:"app_id"`
	URL   string `json:"url"`
}

type LogDrainRes struct {
	Drain LogDrain `json:"drain"`
}

type LogDrainsRes struct {
	Drains []LogDrain `json:"drains"`
}

func (c *Client) LogDrainsList(ctx context.Context, app string) ([]LogDrain, error) {
	var logDrainsRes LogDrainsRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "log_drains", nil, &logDrainsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list the log drains")
	}
	return logDrainsRes.Drains, nil
}

func (c *Client) LogDrainsAddonList(ctx context.Context, app string, addonID string) ([]LogDrain, error) {
	var logDrainsRes LogDrainsRes

	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "addons/"+addonID+"/log_drains", nil, &logDrainsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list the log drains of the addon %s", addonID)
	}
	return logDrainsRes.Drains, nil
}

type LogDrainAddPayload struct {
	Drain LogDrainAddParams `json:"drain"`
}

type LogDrainAddParams struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	Port        string `json:"port"`
	Host        string `json:"host"`
	Token       string `json:"token"`
	DrainRegion string `json:"drain_region"`
}

func (c *Client) LogDrainAdd(ctx context.Context, app string, params LogDrainAddParams) (*LogDrainRes, error) {
	var logDrainRes LogDrainRes
	payload := LogDrainAddPayload{
		Drain: params,
	}

	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "log_drains", payload, &logDrainRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to add drain")
	}

	return &logDrainRes, nil
}

func (c *Client) LogDrainRemove(ctx context.Context, app, URL string) error {
	payload := map[string]string{
		"url": URL,
	}

	req := &httpclient.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/log_drains",
		Expected: httpclient.Statuses{http.StatusNoContent},
		Params:   payload,
	}

	err := c.ScalingoAPI().DoRequest(ctx, req, nil)
	if err != nil {
		return errgo.Notef(err, "fail to delete log drain")
	}

	return nil
}

func (c *Client) LogDrainAddonRemove(ctx context.Context, app, addonID string, URL string) error {
	payload := map[string]string{
		"url": URL,
	}

	req := &httpclient.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/addons/" + addonID + "/log_drains",
		Expected: httpclient.Statuses{http.StatusNoContent},
		Params:   payload,
	}

	err := c.ScalingoAPI().DoRequest(ctx, req, nil)
	if err != nil {
		return errgo.Notef(err, "fail to delete log drain %s from the addon %s", URL, addonID)
	}

	return nil
}

func (c *Client) LogDrainAddonAdd(ctx context.Context, app string, addonID string, params LogDrainAddParams) (*LogDrainRes, error) {
	var logDrainRes LogDrainRes
	payload := LogDrainAddPayload{
		Drain: params,
	}

	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "addons/"+addonID+"/log_drains", payload, &logDrainRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to add log drain to the addon %s", addonID)
	}

	return &logDrainRes, nil
}
