package scalingo

import (
	"time"

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
	req := &APIRequest{
		Client:   c,
		Endpoint: "/features/stacks",
	}

	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	resmap := map[string][]Stack{}
	err = ParseJSON(res, &resmap)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return resmap["stacks"], nil
}
