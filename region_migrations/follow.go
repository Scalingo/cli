package region_migrations

import (
	"github.com/Scalingo/cli/config"
	errgo "gopkg.in/errgo.v1"
)

func Follow(appID, migrationID string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	return WatchMigration(c, appID, migrationID, RefreshOpts{})
}
