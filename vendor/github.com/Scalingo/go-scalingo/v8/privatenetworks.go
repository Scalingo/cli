package scalingo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	httpclient "github.com/Scalingo/go-scalingo/v8/http"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/pagination"
)

type PrivateNetworksService interface {
	PrivateNetworksDomainsList(ctx context.Context, app string, page uint, perPage uint) (pagination.Paginated[[]PrivateNetworkDomain], error)
}

var _ PrivateNetworksService = (*Client)(nil)

type PrivateNetworkDomain = string

type PrivateNetworkDomainsRes struct {
	Domains pagination.Paginated[[]PrivateNetworkDomain] `json:"domain_names"`
}

func (c *Client) PrivateNetworksDomainsList(ctx context.Context, app string, page uint, perPage uint) (pagination.Paginated[[]PrivateNetworkDomain], error) {
	var err error
	validationErr := errors.NewValidationErrorsBuilder()
	if page < 1 {
		validationErr.Set("page", "must be greater than zero")
		return pagination.Paginated[[]PrivateNetworkDomain]{}, validationErr.Build()
	}

	if perPage < 1 || perPage > 50 {
		validationErr.Set("per_page", "must be between 1 and 50")
		return pagination.Paginated[[]PrivateNetworkDomain]{}, validationErr.Build()
	}

	params := url.Values{}
	params.Set("page", strconv.Itoa(int(page)))
	params.Set("per-page", strconv.Itoa(int(perPage)))
	req := &httpclient.APIRequest{
		Method:   http.MethodGet,
		Endpoint: "/apps/" + app + "/private_network_domain_names?" + params.Encode(),
	}
	var domainRes PrivateNetworkDomainsRes
	err = c.ScalingoAPI().DoRequest(ctx, req, &domainRes)
	if err != nil {
		return pagination.Paginated[[]PrivateNetworkDomain]{}, errors.Wrap(ctx, err, "make api call to list the private network domain names")
	}

	return domainRes.Domains, nil
}
