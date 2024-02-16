package alerts

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Remove(ctx context.Context, app, id string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get Scalingo client")
	}

	err = c.AlertRemove(ctx, app, id)
	if err != nil {
		return errors.Wrapf(ctx, err, "remove alert")
	}

	io.Status("The alert has been deleted")
	return nil
}
