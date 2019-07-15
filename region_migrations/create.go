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

	err = WatchMigration(c, app, migration.ID)
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	return nil
}
