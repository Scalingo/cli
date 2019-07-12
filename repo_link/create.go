package repo_link

import (
	"fmt"
	"net/url"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/integrations"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo"
)

func Create(app, integration, link string) error {
	var id string
	var name string

	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		return errgo.New("Source repo url is not a valid http url")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	repoLink, err := c.ScmRepoLinkShow(app)
	if err != nil {
		return errgo.Mask(err)
	}
	if repoLink != nil {
		fmt.Printf("A repo link is already linked with app '%s'.\n", app)
		return nil
	}

	if !utils.IsUUID(integration) {
		i, err := integrations.IntegrationByName(c, integration)
		if err != nil {
			return errgo.Notef(err, "fail to get the integration")
		}

		id = i.ID
		name = i.ScmType
	} else {
		i, err := integrations.IntegrationByUUID(c, integration)
		if err != nil {
			return errgo.Notef(err, "fail to get the integration")
		}

		id = integration
		name = i.ScmType
	}

	_, err = c.ScmRepoLinkAdd(app, scalingo.ScmRepoLinkParams{
		Source:            link,
		AuthIntegrationID: id,
	})
	if err != nil {
		return errgo.Notef(err, "fail to create the repo link")
	}

	fmt.Printf("RepoLink '%s' for app '%s' has been created.\n", name, app)
	return nil
}
