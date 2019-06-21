package domains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func SetCanonical(app, domain string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	d, err := findDomain(c, app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = c.DomainSetCanonical(app, d.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Statusf("Canonical domain set to %s\n", domain)
	return nil
}

func UnsetCanonical(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.DomainUnsetCanonical(app)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Canonical domain disabled")
	return nil
}
