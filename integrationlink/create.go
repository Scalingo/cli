package integrationlink

import (
	"errors"
	"fmt"
	"net/url"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo"
)

func Create(app string, integrationType scalingo.SCMType, integrationURL string, params scalingo.SCMRepoLinkCreateParams) error {
	u, err := url.Parse(integrationURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		return errors.New("source repository URL is not valid")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integration, err := c.SCMIntegrationsShow(string(integrationType))
	if err != nil {
		return errgo.Notef(err, "fail to get the integration")
	}

	repoLink, err := c.SCMRepoLinkShow(app)
	if err != nil {
		return errgo.Notef(err, "fail to get the integration link for this app")
	}
	if repoLink != nil {
		io.Statusf("Your app is already linked to %s/%s/%s", integration.URL, repoLink.Owner, repoLink.Repo)
		if repoLink.Branch != "" {
			fmt.Printf("#%s", repoLink.Branch)
		}
		fmt.Printf(".\n")
		return nil
	}

	params.Source = &integrationURL
	params.AuthIntegrationUUID = &integration.ID

	_, err = c.SCMRepoLinkCreate(app, params)
	if err != nil {
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return errgo.Notef(err, "fail to create the repo link")
		}
		// If error is Payment Required and user tries to exceed its free trial
		return utils.AskAndStopFreeTrial(c, func() error {
			return Create(app, integrationType, integrationURL, params)
		})
	}

	io.Statusf("Your app '%s' is linked to the repository %s.\n", app, integrationURL)
	return nil
}
