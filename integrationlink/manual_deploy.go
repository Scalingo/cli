package integrationlink

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/deployments"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v3"
)

func ManualDeploy(ctx context.Context, app, branch string, follow bool) error {
	if app == "" {
		return errors.New(ctx, "no app defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	deploy, err := c.SCMRepoLinkManualDeploy(ctx, app, branch)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to trigger manual deploy")
	}

	io.Statusf("Manual deployment triggered for app '%s' on branch '%s'.\n", app, branch)

	if !follow {
		return nil
	}

	err = deployments.Stream(ctx, &deployments.StreamOpts{
		AppName:      app,
		DeploymentID: deploy.ID,
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "stream deployment logs")
	}

	return nil
}
