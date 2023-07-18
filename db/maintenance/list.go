package maintenance

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context, app string, addonName string, paginationOpts scalingo.PaginationOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Notef(ctx, err, "get Scalingo client")
	}

	maintenances, pagination, err := c.DatabaseListMaintenance(ctx, app, addonName, paginationOpts)
	if err != nil {
		return errors.Notef(ctx, err, "list the database maintenance")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Type", "Started At", "Ended At", "Status"})

	for _, maintenance := range maintenances {
		startedAt := "Not started"
		if maintenance.StartedAt != nil {
			startedAt = maintenance.StartedAt.Local().Format(utils.TimeFormat)
		}

		endedAt := ""
		if maintenance.EndedAt != nil {
			endedAt = maintenance.EndedAt.Local().Format(utils.TimeFormat)
		}

		t.Append([]string{
			maintenance.ID,
			string(maintenance.Type),
			startedAt,
			endedAt,
			string(maintenance.Status),
		})
	}
	t.Render()
	fmt.Fprintln(os.Stderr, io.Gray(fmt.Sprintf("Page: %d, Last Page: %d", pagination.CurrentPage, pagination.TotalPages)))
	return nil
}
