package scalingo

import (
	"context"
	"errors"
	"time"

	"gopkg.in/errgo.v1"
)

type DomainsService interface {
	DomainsList(ctx context.Context, app string) ([]Domain, error)
	DomainsAdd(ctx context.Context, app string, d Domain) (Domain, error)
	DomainsRemove(ctx context.Context, app string, id string) error
	DomainSetCanonical(ctx context.Context, app, id string) (Domain, error)
	DomainUnsetCanonical(ctx context.Context, app string) (Domain, error)
	DomainSetCertificate(ctx context.Context, app, id, tlsCert, tlsKey string) (Domain, error)
	DomainUnsetCertificate(ctx context.Context, app, id string) (Domain, error)
}

var _ DomainsService = (*Client)(nil)

type LetsEncryptStatus string

const (
	LetsEncryptStatusPendingDNS  LetsEncryptStatus = "pending_dns"
	LetsEncryptStatusNew         LetsEncryptStatus = "new"
	LetsEncryptStatusCreated     LetsEncryptStatus = "created"
	LetsEncryptStatusDNSRequired LetsEncryptStatus = "dns_required"
	LetsEncryptStatusError       LetsEncryptStatus = "error"
)

type SslStatus string

const (
	SslStatusPendingDNS SslStatus = "pending"
	SslStatusNew        SslStatus = "error"
	SslStatusCreated    SslStatus = "success"
)

type ACMEErrorVariables struct {
	DNSProvider string   `json:"dns_provider"`
	Variables   []string `json:"variables"`
}

type Domain struct {
	ID                string             `json:"id"`
	AppID             string             `json:"app_id"`
	Name              string             `json:"name"`
	TLSCert           string             `json:"tlscert,omitempty"`
	TLSKey            string             `json:"tlskey,omitempty"`
	SSL               bool               `json:"ssl"`
	Validity          time.Time          `json:"validity"`
	Canonical         bool               `json:"canonical"`
	LetsEncrypt       bool               `json:"letsencrypt"`
	LetsEncryptStatus LetsEncryptStatus  `json:"letsencrypt_status"`
	SslStatus         SslStatus          `json:"ssl_status"`
	AcmeDNSFqdn       string             `json:"acme_dns_fqdn"`
	AcmeDNSValue      string             `json:"acme_dns_value"`
	AcmeDNSError      ACMEErrorVariables `json:"acme_dns_error"`
}

// domainUnsetCertificateParams is the params of the request to unset the domain certificate.
// We need a dedicated structure rather than using `Domain` so that we remove omitempty from the JSON tags.
type domainUnsetCertificateParams struct {
	TLSCert string `json:"tlscert"`
	TLSKey  string `json:"tlskey"`
}

type DomainsRes struct {
	Domains []Domain `json:"domains"`
}

type DomainRes struct {
	Domain Domain `json:"domain"`
}

func (c *Client) DomainsList(ctx context.Context, app string) ([]Domain, error) {
	var domainRes DomainsRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "domains", nil, &domainRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list the domains")
	}
	return domainRes.Domains, nil
}

func (c *Client) DomainsAdd(ctx context.Context, app string, d Domain) (Domain, error) {
	var domainRes DomainRes
	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "domains", DomainRes{d}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Notef(err, "fail to add a domain")
	}
	return domainRes.Domain, nil
}

func (c *Client) DomainsRemove(ctx context.Context, app, id string) error {
	return c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "domains", id)
}

func (c *Client) DomainsShow(ctx context.Context, app, id string) (Domain, error) {
	var domainRes DomainRes

	err := c.ScalingoAPI().SubresourceGet(ctx, "apps", app, "domains", id, nil, &domainRes)
	if err != nil {
		return Domain{}, errgo.Notef(err, "fail to show the domain")
	}

	return domainRes.Domain, nil
}

func (c *Client) domainsUpdate(ctx context.Context, app, id string, domain Domain) (Domain, error) {
	var domainRes DomainRes
	err := c.ScalingoAPI().SubresourceUpdate(ctx, "apps", app, "domains", id, DomainRes{Domain: domain}, &domainRes)
	if err != nil {
		return Domain{}, errgo.Notef(err, "fail to update the domain")
	}
	return domainRes.Domain, nil
}

func (c *Client) DomainSetCertificate(ctx context.Context, app, id, tlsCert, tlsKey string) (Domain, error) {
	domain, err := c.domainsUpdate(ctx, app, id, Domain{TLSCert: tlsCert, TLSKey: tlsKey})
	if err != nil {
		return Domain{}, errgo.Notef(err, "fail to set the domain certificate")
	}
	return domain, nil
}

func (c *Client) DomainUnsetCertificate(ctx context.Context, app, id string) (Domain, error) {
	var domainRes DomainRes
	err := c.ScalingoAPI().SubresourceUpdate(
		ctx, "apps", app, "domains", id, map[string]domainUnsetCertificateParams{
			"domain": {TLSCert: "", TLSKey: ""},
		}, &domainRes,
	)
	if err != nil {
		return Domain{}, errgo.Notef(err, "fail to unset the domain certificate")
	}
	return domainRes.Domain, nil
}

func (c *Client) DomainSetCanonical(ctx context.Context, app, id string) (Domain, error) {
	domain, err := c.domainsUpdate(ctx, app, id, Domain{Canonical: true})
	if err != nil {
		return Domain{}, errgo.Notef(err, "fail to set the domain as canonical")
	}
	return domain, nil
}

func (c *Client) DomainUnsetCanonical(ctx context.Context, app string) (Domain, error) {
	domains, err := c.DomainsList(ctx, app)
	if err != nil {
		return Domain{}, errgo.Notef(err, "fail to list the domains to unset the canonical one")
	}

	for _, domain := range domains {
		if domain.Canonical {
			domain, err := c.domainsUpdate(ctx, app, domain.ID, Domain{Canonical: false})
			if err != nil {
				return Domain{}, errgo.Notef(err, "fail to unset the domain as canonical")
			}
			return domain, nil
		}
	}
	return Domain{}, errors.New("no canonical domain configured")
}
