package crontasks

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	httpclient "github.com/Scalingo/go-scalingo/v6/http"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context, app string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	cronTasks, err := client.CronTasksGet(ctx, app)
	if err != nil {
		rootError := errors.RootCause(err)
		if !httpclient.IsRequestFailedError(rootError) || rootError.(*httpclient.RequestFailedError).Code != 404 {
			return errors.Wrapf(ctx, err, "fail to get cron tasks")
		}

		// A 404 only means there is no cron task configured on the application. In this case, we want to display an empty table.
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Command", "Size", "Last execution", "Next execution"})

	for _, job := range cronTasks.Jobs {
		lastExecution := job.LastExecutionDate.Format(utils.TimeFormat)
		if job.LastExecutionDate.IsZero() {
			lastExecution = "No previous executions"
		}

		t.Append([]string{job.Command, job.Size, lastExecution, job.NextExecutionDate.Format(utils.TimeFormat)})
	}
	t.Render()

	return nil
}
