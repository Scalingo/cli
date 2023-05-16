package integrationlink

import (
	"context"
	"strconv"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func ManualReviewApp(ctx context.Context, app string, pullRequestID int) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.SCMRepoLinkManualReviewApp(ctx, app, strconv.Itoa(pullRequestID))
	if err != nil {
		return errgo.Notef(err, "fail to manually create a review app")
	}

	io.Statusf("Manual review app created for app '%s' with pull/merge request id '%d'.\n", app, pullRequestID)
	return nil
}
