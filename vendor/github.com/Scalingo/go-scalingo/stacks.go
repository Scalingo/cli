package scalingo

import (
	"time"

	httpclient "github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type Stack struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BaseImage   string    `json:"base_image"`
	Default     bool      `json:"default"`
}

type StacksService interface {
	StacksList() ([]Stack, error)
}

var _ StacksService = (*Client)(nil)

func (c *Client) StacksList() ([]Stack, error) {
	req := &httpclient.APIRequest{
		Endpoint: "/features/stacks",
	}

	resmap := map[string][]Stack{}
	err := c.ScalingoAPI().DoRequest(req, &resmap)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return resmap["stacks"], nil
}
