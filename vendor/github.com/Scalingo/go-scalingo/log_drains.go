package scalingo

import (
	"net/http"

	httpclient "github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type LogDrainsService interface {
	LogDrainsList(app string) ([]LogDrain, error)
	LogDrainAdd(app string, params LogDrainAddParams) (*LogDrainRes, error)
	LogDrainRemove(app, URL string) error
	LogDrainAddonRemove(app, addonID string, URL string) error
	LogDrainsAddonList(app string, addonID string) (LogDrainsRes, error)
	LogDrainAddonAdd(app string, addonID string, params LogDrainAddParams) (*LogDrainRes, error)
}

var _ LogDrainsService = (*Client)(nil)

type LogDrain struct {
	AppID       string `json:"app_id"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Token       string `json:"token"`
	DrainRegion string `json:"drain_region"`
}

type LogDrainRes struct {
	Drain LogDrain `json:"drain"`
}

type LogDrainsRes struct {
	Drains []LogDrain `json:"drains"`
}

func (c *Client) LogDrainsList(app string) ([]LogDrain, error) {
	var logDrainsRes LogDrainsRes
	err := c.ScalingoAPI().SubresourceList("apps", app, "log_drains", nil, &logDrainsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list the log drains")
	}
	return logDrainsRes.Drains, nil
}

func (c *Client) LogDrainsAddonList(app string, addonID string) (LogDrainsRes, error) {
	var logDrainsRes LogDrainsRes

	err := c.ScalingoAPI().SubresourceList("apps", app, "addons/"+addonID+"/log_drains", nil, &logDrainsRes)
	if err != nil {
		return logDrainsRes, errgo.Notef(err, "fail to list the log drains of the addon %s", addonID)
	}
	return logDrainsRes, nil
}

type LogDrainAddParams struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	Port        string `json:"port"`
	Host        string `json:"host"`
	Token       string `json:"token"`
	DrainRegion string `json:"drain_region"`
}

func (c *Client) LogDrainAdd(app string, params LogDrainAddParams) (*LogDrainRes, error) {
	var logDrainRes LogDrainRes
	payload := LogDrainRes{
		Drain: LogDrain{
			Type:        params.Type,
			URL:         params.URL,
			Host:        params.Host,
			Port:        params.Port,
			Token:       params.Token,
			DrainRegion: params.DrainRegion,
		},
	}

	err := c.ScalingoAPI().SubresourceAdd("apps", app, "log_drains", payload, &logDrainRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to add drain")
	}

	return &logDrainRes, nil
}

func (c *Client) LogDrainRemove(app, URL string) error {
	payload := map[string]string{
		"url": URL,
	}

	req := &httpclient.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/log_drains",
		Expected: httpclient.Statuses{http.StatusNoContent},
		Params:   payload,
	}

	err := c.ScalingoAPI().DoRequest(req, nil)
	if err != nil {
		return errgo.Notef(err, "fail to delete log drain")
	}

	return nil
}

func (c *Client) LogDrainAddonRemove(app, addonID string, URL string) error {
	payload := map[string]string{
		"url": URL,
	}

	req := &httpclient.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/addons/" + addonID + "/log_drains",
		Expected: httpclient.Statuses{http.StatusNoContent},
		Params:   payload,
	}

	err := c.ScalingoAPI().DoRequest(req, nil)
	if err != nil {
		return errgo.Notef(err, "fail to delete log drain %s from the addon %s", URL, addonID)
	}

	return nil
}

func (c *Client) LogDrainAddonAdd(app string, addonID string, params LogDrainAddParams) (*LogDrainRes, error) {
	var logDrainRes LogDrainRes
	payload := LogDrainRes{
		Drain: LogDrain{
			Type:        params.Type,
			URL:         params.URL,
			Host:        params.Host,
			Port:        params.Port,
			Token:       params.Token,
			DrainRegion: params.DrainRegion,
		},
	}

	err := c.ScalingoAPI().SubresourceAdd("apps", app, "addons/"+addonID+"/log_drains", payload, &logDrainRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to add log drain to the addon %s", addonID)
	}

	return &logDrainRes, nil
}
