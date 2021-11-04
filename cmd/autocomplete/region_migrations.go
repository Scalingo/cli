package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func RegionMigrationsAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	migrations, err := client.ListRegionMigrations(appName)
	if err != nil {
		return nil
	}

	for _, migration := range migrations {
		fmt.Println(migration.ID)
	}

	return nil

}
