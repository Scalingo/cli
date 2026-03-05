package autoscalers

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	autoscalers, err := c.AutoscalersList(ctx, app)
	if err != nil {
		return errors.Wrapf(ctx, err, "list autoscalers on app %s", app)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Active", "Container type", "Metric", "Target", "Min containers", "Max containers"})

	for _, autoscaler := range autoscalers {
		t.Append([]string{
			strconv.FormatBool(!autoscaler.Disabled),
			autoscaler.ContainerType,
			autoscaler.Metric, fmt.Sprintf("%.2f", autoscaler.Target),
			strconv.Itoa(autoscaler.MinContainers), strconv.Itoa(autoscaler.MaxContainers),
		})
	}
	t.Render()
	return nil
}
