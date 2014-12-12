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

type DomainsRes struct {
	Domains []Domain `json:"urls"`
}

type DomainRes struct {
	Domain Domain `json:"url"`
}

func DomainsList(app string) ([]Domain, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/urls",
		"expected": Statuses{200},
	}
	res, err := Do(req)
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

func DomainsAdd(app string, name string) (Domain, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/urls",
		"expected": Statuses{201},
		"params": map[string]interface{}{
			"url": map[string]interface{}{
				"name": name,
			},
		},
	}

	res, err := Do(req)
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
	req := map[string]interface{}{
		"method":   "DELETE",
		"endpoint": "/apps/" + app + "/urls/" + id,
		"expected": Statuses{204},
	}

	res, err := Do(req)
	if err != nil {
		return errgo.Mask(err)
	}
	res.Body.Close()
	return nil
}
