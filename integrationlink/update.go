package integrationlink

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v7"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Update(ctx context.Context, app string, params scalingo.SCMRepoLinkUpdateParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.SCMRepoLinkUpdate(ctx, app, params)
	if err != nil {
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return errors.Wrapf(ctx, err, "fail to update integration link")
		}

		return utils.AskAndStopFreeTrial(ctx, c, func() error {
			return Update(ctx, app, params)
		})
	}

	io.Statusf("Your app '%s' integration link has been updated.\n", app)
	return nil
}
