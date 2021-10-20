package scalingo

import (
	"fmt"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v4/http"
)

type ContainersService interface {
	ContainersStop(appName, containerID string) error
}

var _ ContainersService = (*Client)(nil)

type Container struct {
	ID        string     `json:"id"`
	AppID     string     `json:"app_id"`
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Command   string     `json:"command"`
	Type      string     `json:"type"`
	TypeIndex int        `json:"type_index"`
	Label     string     `json:"label"`
	State     string     `json:"state"`
	App       *App       `json:"app"`
	// Size has been deprecated in favor of the more comprehensive `ContainerSize` attribute
	Size          string        `json:"size"`
	ContainerSize ContainerSize `json:"container_size"`
}

type ContainerSize struct {
	ID              string `json:"id"`
	HumanCPU        string `json:"human_cpu"`
	Name            string `json:"name"`
	HumanName       string `json:"human_name"`
	Memory          int64  `json:"memory"`
	Ordinal         int    `json:"ordinal"`
	HourlyPrice     int    `json:"hourly_price"`
	ThirtydaysPrice int    `json:"thirtydays_price"`
	PidsLimit       int    `json:"pids_limit,omitempty"`
	Swap            int64  `json:"swap"`
	SKU             string `json:"sku,omitempty"`
}

func (c *Client) ContainersStop(appName, containerID string) error {
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: fmt.Sprintf("/apps/%s/containers/%s/stop", appName, containerID),
		Expected: httpclient.Statuses{202},
	}
	err := c.ScalingoAPI().DoRequest(req, nil)
	if err != nil {
		return errgo.Notef(err, "fail to execute the POST request to stop a container")
	}

	return nil
}
