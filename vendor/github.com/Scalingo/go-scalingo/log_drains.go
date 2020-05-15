package scalingo

import (
	httpclient "github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type LogDrainsService interface {
	LogDrainsList(app string) ([]LogDrain, error)
	LogDrainAdd(app string, params LogDrainAddParams) (*LogDrainRes, error)
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

type logDrainReq struct {
	Drain LogDrain `json:"drain"`
}

type LogDrainRes struct {
	Drain LogDrain `json:"drain"`
	Error string   `json:"error"`
}

func (c *Client) LogDrainsList(app string) ([]LogDrain, error) {
	var logDrainsRes []LogDrain
	err := c.ScalingoAPI().SubresourceList("apps", app, "log_drains", nil, &logDrainsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list the log drains")
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
	payload := logDrainReq{
		Drain: LogDrain{
			Type:        params.Type,
			URL:         params.URL,
			Host:        params.Host,
			Port:        params.Port,
			Token:       params.Token,
			DrainRegion: params.DrainRegion,
		},
	}

	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/log_drains",
		Expected: httpclient.Statuses{201, 422},
		Params:   payload,
	}

	err := c.ScalingoAPI().DoRequest(req, &logDrainRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to add drain")
	}

	if logDrainRes.Error != "" {
		return nil, errgo.Notef(err, logDrainRes.Error)
	}

	return &logDrainRes, nil
}
