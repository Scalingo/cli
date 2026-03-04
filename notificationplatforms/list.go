package notificationplatforms

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	resources, err := c.NotificationPlatformsList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name"})

	for _, r := range resources {
		t.Append([]string{r.Name})
	}
	t.Render()

	return nil
}
