package alerts

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	alerts, err := c.AlertsList(ctx, app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	headers := []string{"ID", "Active", "Container Type", "Metric", "Limit"}
	hasRemindEvery := false
	for _, alert := range alerts {
		if alert.RemindEvery != "" {
			hasRemindEvery = true
		}
	}
	if hasRemindEvery {
		headers = append(headers, "Remind Every")
	}
	t.SetHeader(headers)

	for _, alert := range alerts {
		var above string
		if alert.SendWhenBelow {
			above = "≤"
		} else {
			above = "≥"
		}
		var durationString string
		if alert.DurationBeforeTrigger != 0 {
			durationString = fmt.Sprintf(" (for %s)", alert.DurationBeforeTrigger)
		}

		row := []string{
			alert.ID,
			fmt.Sprint(!alert.Disabled),
			alert.ContainerType,
			alert.Metric,
			fmt.Sprintf("%s %.2f%s", above, alert.Limit, durationString),
		}
		if hasRemindEvery {
			row = append(row, alert.RemindEvery)
		}
		t.Append(row)
	}
	t.Render()
	return nil
}
