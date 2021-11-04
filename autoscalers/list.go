package autoscalers

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func List(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	autoscalers, err := c.AutoscalersList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Active", "Container type", "Metric", "Target", "Min containers", "Max containers"})

	for _, autoscaler := range autoscalers {
		t.Append([]string{
			fmt.Sprint(!autoscaler.Disabled),
			autoscaler.ContainerType,
			autoscaler.Metric, fmt.Sprintf("%.2f", autoscaler.Target),
			fmt.Sprintf("%d", autoscaler.MinContainers), fmt.Sprintf("%d", autoscaler.MaxContainers),
		})
	}
	t.Render()
	return nil
}
