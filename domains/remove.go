package domains

import (
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Remove(app string, domain string) error {
	d, err := findDomain(app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	err = api.DomainsRemove(app, d.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The domain", d.Name, "has been deleted")
	return nil
}

func findDomain(app string, domain string) (*api.Domain, error) {
	domains, err := api.DomainsList(app)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	for _, d := range domains {
		if d.Name == domain {
			return &d, nil
		}
	}
	return nil, errgo.New("no such domain")
}
