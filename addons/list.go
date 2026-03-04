package addons

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
)

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	resources, err := c.AddonsList(ctx, app)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Addon", "ID", "Plan", "Status"})

	for _, resource := range resources {
		t.Append([]string{resource.AddonProvider.Name, resource.ID, resource.Plan.Name, string(resource.Status)})
	}
	t.Render()

	return nil
}
