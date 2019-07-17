package repo_link

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func ManualReviewApp(app, pullRequestID string) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	// Get RepoLink of App
	repoLink, err := c.ScmRepoLinkShow(app)
	if err != nil {
		return errgo.Mask(err)
	}

	err = c.ScmRepoLinkManualReviewApp(app, repoLink.ID, pullRequestID)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Manual review app deployment triggered for app '%s' with pull/merge request id '%s'.\n", app, pullRequestID)
	return nil
}
