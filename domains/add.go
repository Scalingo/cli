package domains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Add(app string, domain string, cert string, key string) error {
	certContent, keyContent, err := validateSSL(cert, key)
	if err != nil {
		return errgo.Mask(err)
	}

	c := config.ScalingoClient()
	d, err := c.DomainsAdd(app, scalingo.Domain{
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
