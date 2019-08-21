package scmintegrations

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"github.com/juju/errgo"
)

func checkIfIntegrationAlreadyExist(c *scalingo.Client, id string) bool {
	integrations, _ := c.SCMIntegrationsShow(id)
	if integrations != nil {
		return true
	}
	return false
}

var (
	// ErrNotFound is returned if no integrations are found
	ErrNotFound = errors.New("not found")
)

func GetTypeFromURL(integrationURL string) (scalingo.SCMType, error) {
	c, err := config.ScalingoClient()
	if err != nil {
		return "", errgo.Notef(err, "fail to get Scalingo client")
	}

	integrations, err := c.SCMIntegrationsList()
	if err != nil {
		return "", errgo.Notef(err, "fail to list SCM integrations")
	}

	u, err := url.Parse(integrationURL)
	if err != nil {
		return "", errgo.Notef(err, "fail to parse the integration URL")
	}
	integrationURLHost := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	for _, i := range integrations {
		if i.URL == integrationURLHost {
			return i.SCMType, nil
		}
	}
	return "", ErrNotFound
}
