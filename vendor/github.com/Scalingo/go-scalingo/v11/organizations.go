package scalingo

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"
)

type OrganizationsService interface {
	OrganizationsList(ctx context.Context) ([]Organization, error)
}

var _ OrganizationsService = (*Client)(nil)

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Role string `json:"role"`
}

type OrganizationsRes struct {
	Organizations []Organization `json:"organizations"`
}

func (c *Client) OrganizationsList(ctx context.Context) ([]Organization, error) {
	var res OrganizationsRes
	err := c.AuthAPI().ResourceList(ctx, "organizations", nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list organizations")
	}
	return res.Organizations, nil
}
