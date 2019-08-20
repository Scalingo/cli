package repo_link

import (
	"net/url"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
)

func Create(app, integrationType, integrationURL string, params scalingo.SCMRepoLinkParams) error {
	u, err := url.Parse(integrationURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		return errgo.New("Source repo url is not a valid http url")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	repoLink, err := c.SCMRepoLinkShow(app)
	if err != nil {
		return errgo.Notef(err, "fail to get repo link for this app")
	}
	if repoLink != nil {
		io.Status("Your app is already linked with an integration.")
		return nil
	}

	integration, err := c.SCMIntegrationsShow(integrationType)
	if err != nil {
		return errgo.Notef(err, "not linked SCM integration or unknown SCM integration")
	}

	params.Source = &integrationURL
	params.AuthIntegrationUUID = &integration.ID

	_, err = c.SCMRepoLinkCreate(app, params)
	if err != nil {
		return errgo.Notef(err, "fail to create the repo link")
	}

	io.Statusf("Repo link with '%s' integration has been created for app '%s'.\n", integration.SCMType, app)
	return nil
}
