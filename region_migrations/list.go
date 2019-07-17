package region_migrations

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	errgo "gopkg.in/errgo.v1"
)

func List(appId string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	migrations, err := c.ListRegionMigrations(appId)
	if err != nil {
		return errgo.Notef(err, "fail to list migrations")
	}
	if len(migrations) == 0 {
		fmt.Println("No migration found for this app")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Destination", "Started At", "Finished At", "Status"})

	for _, migration := range migrations {
		t.Append([]string{
			migration.ID,
			migration.Destination,
			migration.StartedAt.String(),
			migration.FinishedAt.String(),
			formatMigrationStatus(migration.Status),
		})
	}
	t.Render()

	return nil
}
