package scalingo

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/pagination"
)

type PrivateNetworksService interface {
	PrivateNetworksDomainsList(ctx context.Context, app string, paginationReq pagination.Request) (pagination.Paginated[[]PrivateNetworkDomain], error)
}

var _ PrivateNetworksService = (*Client)(nil)

type PrivateNetworkDomain = string

type PrivateNetworkDomainsRes struct {
	Domains pagination.Paginated[[]PrivateNetworkDomain] `json:"domain_names"`
}

func (c *Client) PrivateNetworksDomainsList(ctx context.Context, app string, paginationReq pagination.Request) (pagination.Paginated[[]PrivateNetworkDomain], error) {
	validationErr := errors.NewValidationErrorsBuilder()
	if paginationReq.Page < 1 {
		validationErr.Set("page", "must be greater than zero")
		return pagination.Paginated[[]PrivateNetworkDomain]{}, validationErr.Build()
	}

	if paginationReq.PerPage < 1 || paginationReq.PerPage > 50 {
		validationErr.Set("per_page", "must be between 1 and 50")
		return pagination.Paginated[[]PrivateNetworkDomain]{}, validationErr.Build()
	}

	var domainRes PrivateNetworkDomainsRes
	err := c.ScalingoAPI().SubresourceList(ctx,
		"apps", app, "private_network_domain_names", paginationReq.ToURLValues(),
		&domainRes,
	)
	if err != nil {
		return pagination.Paginated[[]PrivateNetworkDomain]{}, errors.Wrap(ctx, err, "make api call to list the private network domain names")
	}

	return domainRes.Domains, nil
}
