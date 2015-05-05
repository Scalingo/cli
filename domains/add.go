package domains

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
)

func Add(app string, domain string, cert string, key string) error {
	certContent, keyContent, err := validateSSL(cert, key)
	if err != nil {
		return errgo.Mask(err)
	}

	d, err := api.DomainsAdd(app, api.Domain{
		Name:    domain,
		TlsCert: certContent,
		TlsKey:  keyContent,
	})

	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Domain", d.Name, "has been created, access your app at the following URL:\n")
	io.Info("http://" + d.Name)
	return nil
}
