package scalingo

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v5/http"
)

type ContainersService interface {
	ContainersStop(ctx context.Context, appName, containerID string) error
}

var _ ContainersService = (*Client)(nil)

type Container struct {
	ID            string        `json:"id"`
	AppID         string        `json:"app_id"`
	CreatedAt     *time.Time    `json:"created_at"`
	DeletedAt     *time.Time    `json:"deleted_at"`
	Command       string        `json:"command"`
	Type          string        `json:"type"`
	TypeIndex     int           `json:"type_index"`
	Label         string        `json:"label"`
	State         string        `json:"state"`
	App           *App          `json:"app"`
	ContainerSize ContainerSize `json:"container_size"`
}

func (c *Client) ContainersStop(ctx context.Context, appName, containerID string) error {
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: fmt.Sprintf("/apps/%s/containers/%s/stop", appName, containerID),
		Expected: httpclient.Statuses{202},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, nil)
	if err != nil {
		return errgo.Notef(err, "fail to execute the POST request to stop a container")
	}

	return nil
}
