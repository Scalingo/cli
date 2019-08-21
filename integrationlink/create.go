package integrationlink

import (
	"errors"
	"net/url"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
)

func Create(app string, integrationType scalingo.SCMType, integrationURL string, params scalingo.SCMRepoLinkParams) error {
	u, err := url.Parse(integrationURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		return errors.New("source repository URL is not valid")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	repoLink, err := c.SCMRepoLinkShow(app)
	if err != nil {
		return errgo.Notef(err, "fail to get the integration link for this app")
	}
	if repoLink != nil {
		io.Statusf("Your app is already linked to %s/%s#%s.\n", repoLink.Owner, repoLink.Repo, repoLink.Branch)
		return nil
	}

	integration, err := c.SCMIntegrationsShow(string(integrationType))
	if err != nil {
		return errgo.Notef(err, "fail to get the integration")
	}

	params.Source = &integrationURL
	params.AuthIntegrationUUID = &integration.ID

	_, err = c.SCMRepoLinkCreate(app, params)
	if err != nil {
		return errgo.Notef(err, "fail to create the repo link")
	}

	io.Statusf("Your app '%s' is linked to the repository %s.\n", app, integration.URL)
	return nil
}
