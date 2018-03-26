package alerts

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List(app string) error {
	c := config.ScalingoClient()
	alerts, err := c.AlertsList(app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	headers := []string{"ID", "Active", "Container Type", "Metric", "Limit"}
	hasRemindEvery := false
	for _, alert := range alerts {
		if alert.RemindEvery != "" {
			headers = append(headers, "Remind Every")
			hasRemindEvery = true
		}
	}
	t.SetHeader(headers)

	for _, alert := range alerts {
		var above string
		if alert.SendWhenBelow {
			above = "below"
		} else {
			above = "above"
		}
		row := []string{
			alert.ID,
			fmt.Sprint(!alert.Disabled),
			alert.ContainerType,
			alert.Metric,
			fmt.Sprintf("triggers %s %.2f", above, alert.Limit),
		}
		if hasRemindEvery {
			row = append(row, alert.RemindEvery)
		}
		t.Append(row)
	}
	t.Render()
	return nil
}
