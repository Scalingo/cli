package scalingo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v8/http"
	"github.com/Scalingo/go-utils/pagination"
)

type PrivateNetworksService interface {
	PrivateNetworksDomainsList(ctx context.Context, app string, page string, perPage string) (pagination.Paginated[[]PrivateNetworkDomain], error)
}

var _ PrivateNetworksService = (*Client)(nil)

type PrivateNetworkDomain = string

type PrivateNetworkDomainsRes struct {
	Domains pagination.Paginated[[]PrivateNetworkDomain] `json:"domain_names"`
}

func (c *Client) PrivateNetworksDomainsList(ctx context.Context, app string, page string, perPage string) (pagination.Paginated[[]PrivateNetworkDomain], error) {
	var err error
	pageInt := 1
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			return pagination.Paginated[[]PrivateNetworkDomain]{}, errgo.Newf("invalid page number: %s", page)
		}
	}

	if perPage != "" {
		perPageInt, err := strconv.Atoi(perPage)
		if err != nil || perPageInt < 1 || perPageInt > 50 {
			return pagination.Paginated[[]PrivateNetworkDomain]{}, errgo.Newf("invalid per_page number: %s", perPage)
		}
	}

	params := url.Values{}
	params.Set("page", strconv.Itoa(pageInt))
	params.Set("per-page", perPage)
	req := &httpclient.APIRequest{
		Method:   http.MethodGet,
		Endpoint: "/apps/" + app + "/private-network-domain-names?" + params.Encode(),
	}
	var domainRes PrivateNetworkDomainsRes
	err = c.ScalingoAPI().DoRequest(ctx, req, &domainRes)
	if err != nil {
		return pagination.Paginated[[]PrivateNetworkDomain]{}, errgo.Mask(err)
	}

	return domainRes.Domains, nil
}
