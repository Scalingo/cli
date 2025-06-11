package alerts

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	alerts, err := c.AlertsList(ctx, app)
	if err != nil {
		return errors.Wrap(ctx, err, "list alerts")
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
	t.Header(headers)

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
			strconv.FormatBool(!alert.Disabled),
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
