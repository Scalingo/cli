package integrationlink

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v4"
)

func Update(app string, params scalingo.SCMRepoLinkUpdateParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.SCMRepoLinkUpdate(app, params)
	if err != nil {
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return errgo.Notef(err, "fail to update integration link")
		}

		return utils.AskAndStopFreeTrial(c, func() error {
			return Update(app, params)
		})
	}

	io.Statusf("Your app '%s' integration link has been updated.\n", app)
	return nil
}
