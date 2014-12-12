package domains

import (
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Add(app string, domain string) error {
	d, err := api.DomainsAdd(app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Domain", d.Name, "has been created, access your app at the following URL:\n")
	io.Info("http://" + d.Name)
	return nil
}
