package region_migrations

import (
	"context"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Follow(ctx context.Context, appID, migrationID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	return WatchMigration(ctx, c, appID, migrationID, RefreshOpts{
		ShowHints: true,
	})
}
