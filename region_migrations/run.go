package region_migrations

import (
	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo"
	errgo "gopkg.in/errgo.v1"
)

func Create(app string, destination string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	migration, err := c.CreateRegionMigration(app, scalingo.RegionMigrationParams{
		Destination: destination,
	})
	if err != nil {
		return errgo.Notef(err, "fail to create migration")
	}

	err = c.RunRegionMigrationStep(app, migration.ID, scalingo.RegionMigrationStepPreflight)
	if err != nil {
		return errgo.Notef(err, "fail to run preflight step")
	}

	err = WatchMigration(c, app, migration.ID, RefreshOpts{
		ExpectedStatuses: []scalingo.RegionMigrationStatus{
			scalingo.RegionMigrationStatusPreflightError,
			scalingo.RegionMigrationStatusPreflightSuccess,
		},
	})
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	return nil
}

func Run(app, migrationId string, step scalingo.RegionMigrationStep) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	err = c.RunRegionMigrationStep(app, migrationId, step)
	if err != nil {
		return errgo.Notef(err, "fail to run %s step", step)
	}

	err = WatchMigration(c, app, migrationId, RefreshOpts{
		ExpectedStatuses: []scalingo.RegionMigrationStatus{},
	})
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	return nil
}

func Abort(app, migrationId string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	err = c.RunRegionMigrationStep(app, migrationId, scalingo.RegionMigrationStepAbort)
	if err != nil {
		return errgo.Notef(err, "fail to run abort step")
	}

	err = WatchMigration(c, app, migrationId, RefreshOpts{
		ExpectedStatuses: []scalingo.RegionMigrationStatus{},
	})
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	return nil

}
