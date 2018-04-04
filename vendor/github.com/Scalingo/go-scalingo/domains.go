package scalingo

import (
	"time"

	"gopkg.in/errgo.v1"
)

type DomainsService interface {
	DomainsList(app string) ([]Domain, error)
	DomainsAdd(app string, d Domain) (Domain, error)
	DomainsRemove(app string, id string) error
	DomainsUpdate(app, id, cert, key string) (Domain, error)
}

var _ DomainsService = (*Client)(nil)

type Domain struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	TlsCert  string    `json:"tlscert,omitempty"`
	TlsKey   string    `json:"tlskey,omitempty"`
	SSL      bool      `json:"ssl"`
	Validity time.Time `json:"validity"`
}

type DomainsRes struct {
	Domains []Domain `json:"domains"`
}

type DomainRes struct {
	Domain Domain `json:"domain"`
}

func (c *Client) DomainsList(app string) ([]Domain, error) {
	var domainRes DomainsRes
	err := c.subresourceList(app, "domains", nil, &domainRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return domainRes.Domains, nil
}

func (c *Client) DomainsAdd(app string, d Domain) (Domain, error) {
	var domainRes DomainRes
	err := c.subresourceAdd(app, "domains", DomainRes{d}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}
	return domainRes.Domain, nil
}

func (c *Client) DomainsRemove(app string, id string) error {
	return c.subresourceDelete(app, "domains", id)
}

func (c *Client) DomainsUpdate(app, id, cert, key string) (Domain, error) {
	var domainRes DomainRes
	err := c.subresourceUpdate(app, "domains", id, DomainRes{Domain: Domain{TlsCert: cert, TlsKey: key}}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}
	return domainRes.Domain, nil
}
