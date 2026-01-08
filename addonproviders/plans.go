package addonproviders

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Plans(ctx context.Context, addon string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	plans, err := c.AddonProviderPlansList(ctx, addon, scalingo.AddonProviderPlansListOpts{})
	if err != nil {
		return errors.Wrapf(ctx, err, "list addon provider plans")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Name"})
	for _, plan := range plans {
		t.Append([]string{plan.Name, plan.DisplayName})
	}
	t.Render()
	return nil
}
