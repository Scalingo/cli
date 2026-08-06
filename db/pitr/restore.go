package pitr

import (
	"context"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Restore(ctx context.Context, currentResource, addonName string, restoreTime time.Time) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get Scalingo client")
	}

	operationID, err := c.DatabaseRestorePITR(ctx, currentResource, addonName, restoreTime)
	if err != nil {
		return errors.Wrap(ctx, err, "restore database")
	}

	io.Statusf("Database restore operation %s has been created\n", operationID)

	return nil
}
