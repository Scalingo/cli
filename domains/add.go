package domains

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v5"
)

func Add(ctx context.Context, app string, domain string, cert string, key string) error {
	certContent, keyContent, err := validateSSL(cert, key)
	if err != nil {
		return errgo.Mask(err)
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	d, err := c.DomainsAdd(ctx, app, scalingo.Domain{
		Name:    domain,
		TLSCert: certContent,
		TLSKey:  keyContent,
	})

	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Domain", d.Name, "has been created, access your app at the following URL:\n")
	io.Info("http://" + d.Name)
	return nil
}
