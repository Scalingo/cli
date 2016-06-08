package domains

import (
	"io/ioutil"

	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func DisableSSL(app string, domain string) error {
	d, err := findDomain(app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	c := config.ScalingoClient()
	_, err = c.DomainsUpdate(app, d.ID, "", "")
	if err != nil {
		return errgo.Mask(err)
	}
	io.Status("SSL of " + domain + " has been disabled.")
	return nil
}

func EnableSSL(app, domain, certPath, keyPath string) error {
	d, err := findDomain(app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	certContent, keyContent, err := validateSSL(certPath, keyPath)
	if err != nil {
		return errgo.Mask(err)
	}

	c := config.ScalingoClient()
	d, err = c.DomainsUpdate(app, d.ID, certContent, keyContent)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The certificate and key have been installed for " + d.Name + " (Validity: " + d.Validity.UTC().String() + ")")
	return nil
}

func validateSSL(cert, key string) (string, string, error) {
	if cert == "" && key == "" {
		return "", "", nil
	}

	if cert == "" && key != "" {
		return "", "", errgo.New("--cert <certificate path> should be defined")
	}

	if key == "" && cert != "" {
		return "", "", errgo.New("--key <key path> should be defined")
	}

	certContent, err := ioutil.ReadFile(cert)
	if err != nil {
		return "", "", errgo.Mask(err)
	}
	keyContent, err := ioutil.ReadFile(key)
	if err != nil {
		return "", "", errgo.Mask(err)
	}
	return string(certContent), string(keyContent), nil
}
