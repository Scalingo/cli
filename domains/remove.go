package domains

import (
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Remove(app string, domain string) error {
	domains, err := api.DomainsList(app)
	if err != nil {
		return errgo.Mask(err)
	}

	d := findDomain(domains, domain)
	if d == nil {
		return errgo.New("no such domain")
	}

	err = api.DomainsRemove(app, d.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The domain", d.Name, "has been deleted")
	return nil
}

func findDomain(domains []api.Domain, domain string) *api.Domain {
	for _, d := range domains {
		if d.Name == domain {
			return &d
		}
	}
	return nil
}
