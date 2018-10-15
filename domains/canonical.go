package domains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func SetCanonical(app, domain string) error {
	d, err := findDomain(app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	c := config.ScalingoClient()
	_, err = c.DomainSetCanonical(app, d.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Statusf("Canonical domain set to %s\n", domain)
	return nil
}

func UnsetCanonical(app string) error {
	c := config.ScalingoClient()

	_, err := c.DomainUnsetCanonical(app)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Canonical domain disabled")
	return nil
}
