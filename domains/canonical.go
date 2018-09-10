package domains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Canonical(app, domain string, enable bool) error {
	d, err := findDomain(app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	c := config.ScalingoClient()
	if enable {
		_, err = c.DomainSetCanonical(app, d.ID)
	} else {
		_, err = c.DomainUnsetCanonical(app, d.ID)
	}
	if err != nil {
		return errgo.Mask(err)
	}

	if enable {
		io.Statusf("Canonical domain set to %s", domain)
	} else {
		io.Status("Canonical domain disabled")
	}
	return nil
}
