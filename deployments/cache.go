package deployments

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func ResetCache(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	err = c.DeploymentCacheReset(ctx, app)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
