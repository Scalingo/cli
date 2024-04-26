package domains

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func DisableSSL(ctx context.Context, app string, domain string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	d, err := findDomain(ctx, c, app, domain)
	if err != nil {
		return errgo.Notef(err, "fail to find the matching domain to disable SSL")
	}

	_, err = c.DomainUnsetCertificate(ctx, app, d.ID)
	if err != nil {
		return errgo.Notef(err, "fail to unset the domain certificate")
	}
	io.Status("SSL of " + domain + " has been disabled.")
	return nil
}

func EnableSSL(ctx context.Context, app, domain, cert, key string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	d, err := findDomain(ctx, c, app, domain)
	if err != nil {
		return errgo.Notef(err, "fail to find the matching domain to enable SSL")
	}

	d, err = c.DomainSetCertificate(ctx, app, d.ID, cert, key)
	if err != nil {
		return errgo.Notef(err, "fail to set the domain certificate")
	}

	io.Status("The certificate and key have been installed for " + d.Name + " (Validity: " + d.Validity.UTC().String() + ")")
	return nil
}
