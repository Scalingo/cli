package api

import (
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

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

func DomainsList(app string) ([]Domain, error) {
	var domainRes DomainsRes
	err := subresourceList(app, "domains", nil, &domainRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return domainRes.Domains, nil
}

func DomainsAdd(app string, d Domain) (Domain, error) {
	var domainRes DomainRes
	err := subresourceAdd(app, "domains", DomainRes{d}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}
	return domainRes.Domain, nil
}

func DomainsRemove(app string, id string) error {
	return subresourceDelete(app, "domains", id)
}

func DomainsUpdate(app, id, cert, key string) (Domain, error) {
	var domainRes DomainRes
	err := subresourceUpdate(app, "domains", id, DomainRes{Domain: Domain{TlsCert: cert, TlsKey: key}}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}
	return domainRes.Domain, nil
}
