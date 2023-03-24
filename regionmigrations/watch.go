package regionmigrations

import (
	"context"
	"fmt"

	errgo "gopkg.in/errgo.v1"

	scalingo "github.com/Scalingo/go-scalingo/v6"
)

func WatchMigration(ctx context.Context, client *scalingo.Client, appID, migrationID string, opts RefreshOpts) error {
	refresher := NewRefresher(client, appID, migrationID, opts)
	migration, err := refresher.Start()
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	if migration == nil {
		return nil
	}

	migrationFinished(ctx, appID, *migration, opts)

	return nil
}

func migrationFinished(ctx context.Context, appID string, migration scalingo.RegionMigration, opts RefreshOpts) {
	fmt.Printf("\n\n")
	switch migration.Status {
	case scalingo.RegionMigrationStatusDone:
		showMigrationStatusSuccess(ctx, appID, migration)
	case scalingo.RegionMigrationStatusError:
		fallthrough
	case scalingo.RegionMigrationStatusPreflightError:
		showMigrationStatusFailed(appID, migration, opts)
	case scalingo.RegionMigrationStatusPreflightSuccess:
		showMigrationStatusPreflightSuccess(appID, migration)
	case scalingo.RegionMigrationStatusPrepared:
		showMigrationStatusPrepared(appID, migration)
	case scalingo.RegionMigrationStatusDataMigrated:
		showMigrationStatusDataMigrated(appID, migration)
	case scalingo.RegionMigrationStatusAborted:
		showMigrationStatusAborted(appID, migration)
	}
}
