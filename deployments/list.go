package deployments

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context, app string, paginationOpts scalingo.PaginationOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	deployments, pagination, err := c.DeploymentListWithPagination(ctx, app, paginationOpts)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list the application deployments")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Date", "Duration", "User", "Git Ref", "Status", "Image Size"})

	for _, deployment := range deployments {
		var duration string
		if deployment.Duration != 0 {
			d, _ := time.ParseDuration(fmt.Sprintf("%ds", deployment.Duration))
			duration = d.String()
		} else {
			duration = "n/a"
		}
		t.Append([]string{deployment.ID,
			deployment.CreatedAt.Format(utils.TimeFormat),
			duration,
			deployment.User.Username,
			deployment.GitRef,
			string(deployment.Status),
			humanize.IBytes(deployment.ImageSize),
		})
	}
	t.Render()
	fmt.Fprintln(os.Stderr, io.Gray(fmt.Sprintf("Page: %d, Last Page: %d", pagination.CurrentPage, pagination.TotalPages)))
	return nil
}
