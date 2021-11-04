package region_migrations

import (
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Follow(appID, migrationID string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	return WatchMigration(c, appID, migrationID, RefreshOpts{
		ShowHints: true,
	})
}
