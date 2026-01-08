package integrationlink

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

func PullRequest(ctx context.Context, app string, pullRequestID int) (*scalingo.RepoLinkPullRequest, error) {
	if app == "" {
		return nil, errgo.New("no app defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	pullRequest, err := c.SCMRepoLinkPullRequest(ctx, app, pullRequestID)

	if err != nil {
		return nil, errors.Wrapf(ctx, err, "fail to fetch the pull request status")
	}

	return pullRequest, nil
}
