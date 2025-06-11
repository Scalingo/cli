package notificationplatforms

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	resources, err := c.NotificationPlatformsList(ctx)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name"})

	for _, r := range resources {
		t.Append([]string{r.Name})
	}
	t.Render()

	return nil
}
