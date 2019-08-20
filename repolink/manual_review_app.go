package repolink

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func ManualReviewApp(app, pullRequestID string) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.SCMRepoLinkManualReviewApp(app, pullRequestID)
	if err != nil {
		return errgo.Notef(err, "fail to manually create a review app")
	}

	io.Statusf("Manual review app created for app '%s' with pull/merge request id '%s'.\n", app, pullRequestID)
	return nil
}
