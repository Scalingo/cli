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
	t.SetHeader([]string{"ID", "Active", "Container Type", "Metric", "Limit"})

	for _, alert := range alerts {
		var above string
		if alert.SendWhenBelow {
			above = "below"
		} else {
			above = "above"
		}
		t.Append([]string{
			alert.ID,
			fmt.Sprint(!alert.Disabled),
			alert.ContainerType,
			alert.Metric,
			fmt.Sprintf("triggers %s %.2f", above, alert.Limit),
		})
	}
	t.Render()
	return nil
}
