package integrationlink

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Delete(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	err = c.SCMRepoLinkDelete(ctx, app)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to delete repo link")
	}

	io.Statusf("Current integration link has been deleted from app '%s'.\n", app)
	return nil
}
