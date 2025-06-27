package scmintegrations

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
)

var ErrNotFound = errors.New("SCM integration not found")

func checkIfIntegrationAlreadyExist(ctx context.Context, c *scalingo.Client, id string) bool {
	integrations, _ := c.SCMIntegrationsShow(ctx, id)
	return integrations != nil
}

func GetTypeFromURL(ctx context.Context, integrationURL string) (scalingo.SCMType, error) {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return "", errgo.Notef(err, "fail to get Scalingo client")
	}

	integrations, err := c.SCMIntegrationsList(ctx)
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
