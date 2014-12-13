package domains

import (
	"io/ioutil"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func DisableSSL(app string, domain string) error {
	d, err := findDomain(app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = api.DomainsUpdate(app, d.ID, "", "")
	if err != nil {
		return errgo.Mask(err)
	}
	io.Status("SSL of " + domain + " has been disabled.")
	return nil
}

func EnableSSL(app, domain, certPath, keyPath string) error {
	d, err := findDomain(app, domain)
	certContent, err := ioutil.ReadFile(certPath)
	if err != nil {
		return errgo.Mask(err)
	}
	keyContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return errgo.Mask(err)
	}
	_, err = api.DomainsUpdate(app, d.ID, string(certContent), string(keyContent))
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
