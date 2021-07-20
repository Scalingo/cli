package crontasks

import (
	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"os"
)

func List(app string) error {
	client, err := config.ScalingoClient()
	if err != nil {
		return errors.Wrap(err, "fail to get Scalingo client")
	}

	cronTasks, err := client.CronTasksGet(app)
	if err != nil {
		return errors.Wrap(err, "fail to get cron tasks")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Command", "Size"})

	for _, job := range cronTasks.Jobs {
		t.Append([]string{job.Command, job.Size})
	}
	t.Render()

	return nil
}
