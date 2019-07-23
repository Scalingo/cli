package region_migrations

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/fatih/color"
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
		io.Status("No migration found for this app")
		return nil
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].StartedAt.Unix() > migrations[j].StartedAt.Unix()
	})

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Destination", "Started At", "Finished At", "Status"})

	for _, migration := range migrations {
		finishedAt := migration.FinishedAt.Local().Format(time.RFC822)
		if migration.FinishedAt.IsZero() {
			finishedAt = color.BlueString("Ongoing")
		}

		t.Append([]string{
			migration.ID,
			migration.Destination,
			migration.StartedAt.Local().Format(time.RFC822),
			finishedAt,
			formatMigrationStatus(migration.Status),
		})
	}
	t.Render()

	return nil
}
