package pitr

import (
	"context"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func GetRecoveryWindow(ctx context.Context, currentResource, addonName string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get Scalingo client")
	}

	recoveryWindow, err := c.DatabaseGetPITRRecoveryWindow(ctx, currentResource, addonName)
	if err != nil {
		return errors.Wrap(ctx, err, "get PITR recovery window")
	}

	var earliestRecoverableAt string
	if recoveryWindow.EarliestRecoverableAt != nil {
		earliestRecoverableAt = recoveryWindow.EarliestRecoverableAt.Format(time.RFC3339)
	}
	var latestRecoverableAt string
	if recoveryWindow.LatestRecoverableAt != nil {
		latestRecoverableAt = recoveryWindow.LatestRecoverableAt.Format(time.RFC3339)
	}

	t := tablewriter.NewWriter(os.Stdout)
	_ = t.Append([]string{"Earliest recoverable at", earliestRecoverableAt})
	_ = t.Append([]string{"Latest recoverable at", latestRecoverableAt})

	_ = t.Render()

	return nil
}
