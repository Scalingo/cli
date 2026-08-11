package pitr

import (
	"context"
	"fmt"
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

	io.Warning("This operation is irreversible and previous data will be restore")
	io.Info("Do you want to confirm? (y/n)")

	var confirm string
	_, _ = fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		return errors.New(ctx, "You didn't confirm, aborting…")
	}

	operationID, err := c.DatabaseRestorePITR(ctx, currentResource, addonName, restoreTime)
	if err != nil {
		return errors.Wrap(ctx, err, "restore database")
	}

	io.Statusf("Database restore operation %s has been created\n", operationID)

	return nil
}
