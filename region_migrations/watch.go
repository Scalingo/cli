package region_migrations

import (
	"fmt"

	scalingo "github.com/Scalingo/go-scalingo/v4"
	errgo "gopkg.in/errgo.v1"
)

func WatchMigration(client *scalingo.Client, appId, migrationId string, opts RefreshOpts) error {
	refresher := NewRefresher(client, appId, migrationId, opts)
	migration, err := refresher.Start()
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	if migration == nil {
		return nil
	}

	migrationFinished(appId, *migration, opts)

	return nil
}

func migrationFinished(appId string, migration scalingo.RegionMigration, opts RefreshOpts) {
	fmt.Printf("\n\n")
	switch migration.Status {
	case scalingo.RegionMigrationStatusDone:
		showMigrationStatusSuccess(appId, migration)
	case scalingo.RegionMigrationStatusError:
		fallthrough
	case scalingo.RegionMigrationStatusPreflightError:
		showMigrationStatusFailed(appId, migration, opts)
	case scalingo.RegionMigrationStatusPreflightSuccess:
		showMigrationStatusPreflightSuccess(appId, migration)
	case scalingo.RegionMigrationStatusPrepared:
		showMigrationStatusPrepared(appId, migration)
	case scalingo.RegionMigrationStatusDataMigrated:
		showMigrationStatusDataMigrated(appId, migration)
	case scalingo.RegionMigrationStatusAborted:
		showMigrationStatusAborted(appId, migration)
	}
}
