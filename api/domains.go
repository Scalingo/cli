package api

import (
	"time"

	"gopkg.in/errgo.v1"
)

type Domain struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	SSL      bool      `json:"ssl"`
	Validity time.Time `json:"validity"`
}

type DomainCreate struct {
	Name    string `json:"name"`
	TlsCert string `json:"tlscert,omitempty"`
	TlsKey  string `json:"tlskey,omitempty"`
}

type DomainsRes struct {
	Domains []Domain `json:"domains"`
}

type DomainRes struct {
	Domain Domain `json:"domain"`
}

func DomainsList(app string) ([]Domain, error) {
	req := &APIRequest{
		Endpoint: "/apps/" + app + "/domains",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var domainRes DomainsRes
	err = ParseJSON(res, &domainRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return domainRes.Domains, nil
}

func DomainsAdd(app string, d *DomainCreate) (Domain, error) {
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/domains",
		Expected: Statuses{201},
		Params: map[string]interface{}{
			"domain": d,
		},
	}

	res, err := req.Do()
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}
	defer res.Body.Close()

	var domainRes DomainRes
	err = ParseJSON(res, &domainRes)
	if err != nil {
		return Domain{}, errgo.Mask(err)
	}

	return domainRes.Domain, nil
}

func DomainsRemove(app string, id string) error {
	req := &APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/domains/" + id,
		Expected: Statuses{204},
	}

	res, err := req.Do()
	if err != nil {
		return errgo.Mask(err)
	}
	res.Body.Close()
	return nil
}

func DomainsUpdate(app, id, cert, key string) (*Domain, error) {
	req := &APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app + "/domains/" + id,
		Params: map[string]interface{}{
			"domain": map[string]interface{}{
				"tlscert": cert,
				"tlskey":  key,
			},
		},
	}

	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	var domainRes DomainRes
	err = ParseJSON(res, &domainRes)
	if err != nil {
		return nil, nil
	}

	return &domainRes.Domain, nil
}
