package scmintegrations

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func checkIfIntegrationAlreadyExist(c *scalingo.Client, id string) bool {
	integrations, _ := c.SCMIntegrationsShow(id)
	if integrations != nil {
		return true
	}
	return false
}

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

	// If no integration matches the given URL, display a helpful status
	// message
	switch u.Host {
	case "github.com":
		io.Error("No GitHub integration found, please follow this URL to add it:")
		io.Errorf("%s/users/github/link\n", config.C.ScalingoAuthUrl)
	case "gitlab.com":
		io.Error("No GitLab integration found, please follow this URL to add it:")
		io.Errorf("%s/users/gitlab/link\n", config.C.ScalingoAuthUrl)
	default:
		io.Errorf("No integration found for URL %s.\n", integrationURL)
		io.Errorf("Please run the command:\n")
		io.Errorf("scalingo integrations-add gitlab-self-hosted|github-enterprise --url %s://%s --token <personal-access-token>\n", u.Scheme, u.Host)
	}
	return "", errors.New("not found")
}
