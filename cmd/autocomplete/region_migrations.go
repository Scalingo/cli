package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v2"
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func RegionMigrationsAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	migrations, err := client.ListRegionMigrations(c.Context, appName)
	if err != nil {
		return nil
	}

	for _, migration := range migrations {
		fmt.Println(migration.ID)
	}

	return nil
}
