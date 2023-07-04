package maintenance

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
)

func List(ctx context.Context, app string, addonName string, paginationOpts scalingo.PaginationOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	response, err := c.DatabaseListMaintenance(ctx, app, addonName, paginationOpts)
	if err != nil {
		return errgo.Notef(err, "fail to list the database maintenance")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Type", "Started At", "Ended At", "Status"})

	for _, maintenance := range response.Maintenance {
		startedAt := "Not started"
		if maintenance.StartedAt != nil {
			startedAt = maintenance.StartedAt.Format(utils.TimeFormat)
		}

		endedAt := ""
		if maintenance.EndedAt != nil {
			endedAt = maintenance.EndedAt.Format(utils.TimeFormat)
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
	fmt.Fprintln(os.Stderr, io.Gray(fmt.Sprintf("Page: %d, Last Page: %d", response.Meta.CurrentPage, response.Meta.TotalPages)))
	return nil
}
