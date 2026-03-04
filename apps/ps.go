package apps

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
)

func Ps(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client to list the application containers")
	}

	containers, err := c.AppsContainersPs(ctx, app)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list the application containers")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Status", "Command", "Size", "Created At"})

	for _, container := range containers {
		t.Append([]string{container.Label, container.State, container.Command, container.ContainerSize.HumanName, container.CreatedAt.Format(utils.TimeFormat)})
	}
	t.Render()
	return nil
}
