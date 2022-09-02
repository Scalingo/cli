package integrationlink

import (
	"context"
	"fmt"
	"net/url"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v5"
	"github.com/Scalingo/go-scalingo/v5/http"
	"github.com/Scalingo/go-utils/errors"
)

func Create(ctx context.Context, app string, integrationType scalingo.SCMType, integrationURL string, params scalingo.SCMRepoLinkCreateParams) error {
	u, err := url.Parse(integrationURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		return errgo.New("source repository URL is not valid")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integration, err := c.SCMIntegrationsShow(ctx, string(integrationType))
	if err != nil {
		return errgo.Notef(err, "fail to get the integration")
	}

	repoLink, err := c.SCMRepoLinkShow(ctx, app)
	if err != nil {
		scerr, ok := errors.ErrgoRoot(err).(*http.RequestFailedError)
		if !ok || scerr.Code != 404 {
			return errgo.Notef(err, "fail to get the integration link for this app")
		}
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

	_, err = c.SCMRepoLinkCreate(ctx, app, params)
	if err != nil {
		if utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return utils.AskAndStopFreeTrial(ctx, c, func() error {
				return Create(ctx, app, integrationType, integrationURL, params)
			})
		}

		return errgo.Notef(err, "fail to create the repo link")
	}

	io.Statusf("Your app '%s' is linked to the repository %s.\n", app, integrationURL)
	return nil
}
