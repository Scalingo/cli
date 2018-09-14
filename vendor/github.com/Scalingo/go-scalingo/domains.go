package scalingo

import (
	"errors"
	"time"

	"gopkg.in/errgo.v1"
)

type DomainsService interface {
	DomainsList(app string) ([]Domain, error)
	DomainsAdd(app string, d Domain) (Domain, error)
	DomainsRemove(app string, id string) error
	DomainsUpdate(app, id, cert, key string) (Domain, error)
	DomainSetCanonical(app, id string) (Domain, error)
	DomainUnsetCanonical(app string) (Domain, error)
}

var _ DomainsService = (*Client)(nil)

type Domain struct {
	ID        string    `json:"id"`
	AppID     string    `json:"app_id"`
	Name      string    `json:"name"`
	TLSCert   string    `json:"tlscert,omitempty"`
	TLSKey    string    `json:"tlskey,omitempty"`
	SSL       bool      `json:"ssl"`
	Validity  time.Time `json:"validity"`
	Canonical bool      `json:"canonical"`
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

func (c *Client) DomainsRemove(app, id string) error {
	return c.subresourceDelete(app, "domains", id)
}

func (c *Client) DomainsUpdate(app, id, cert, key string) (Domain, error) {
	var domainRes DomainRes
	err := c.subresourceUpdate(app, "domains", id, DomainRes{Domain: Domain{TLSCert: cert, TLSKey: key}}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}
	return domainRes.Domain, nil
}

func (c *Client) DomainsShow(app, id string) (Domain, error) {
	var domainRes DomainRes

	err := c.subresourceGet(app, "domains", id, nil, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}

	return domainRes.Domain, nil
}

func (c *Client) DomainSetCanonical(app, id string) (Domain, error) {
	var domainRes DomainRes
	err := c.subresourceUpdate(app, "domains", id, DomainRes{Domain: Domain{Canonical: true}}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}
	return domainRes.Domain, nil
}

func (c *Client) DomainUnsetCanonical(app string) (Domain, error) {
	domains, err := c.DomainsList(app)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}

	for _, domain := range domains {
		if domain.Canonical {
			var domainRes DomainRes
			err := c.subresourceUpdate(app, "domains", domain.ID, DomainRes{Domain: Domain{Canonical: false}}, &domainRes)
			if err != nil {
				return Domain{}, errgo.Mask(err)
			}
			return domainRes.Domain, nil
		}
	}
	return Domain{}, errors.New("no canonical domain configured")
}
