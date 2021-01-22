package domains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func Remove(app string, domain string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	d, err := findDomain(c, app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	err = c.DomainsRemove(app, d.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The domain", d.Name, "has been deleted")
	return nil
}

func findDomain(c *scalingo.Client, app string, domain string) (scalingo.Domain, error) {
	domains, err := c.DomainsList(app)
	if err != nil {
		return scalingo.Domain{}, errgo.Mask(err)
	}

	for _, d := range domains {
		if d.Name == domain {
			return d, nil
		}
	}
	return scalingo.Domain{}, errgo.New("There is no such domain, please ensure you've added it correctly.\nhttps://my.scalingo.com/apps/" + app + "/domains")
}
