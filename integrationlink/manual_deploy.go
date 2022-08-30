package integrationlink

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func ManualDeploy(ctx context.Context, app, branch string) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.SCMRepoLinkManualDeploy(ctx, app, branch)
	if err != nil {
		return errgo.Notef(err, "fail to trigger manual deploy")
	}

	io.Statusf("Manual deployment triggered for app '%s' on branch '%s'.\n", app, branch)
	return nil
}
