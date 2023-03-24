package addonproviders

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Plans(ctx context.Context, addon string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	plans, err := c.AddonProviderPlansList(ctx, addon)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name"})
	for _, plan := range plans {
		t.Append([]string{plan.Name, plan.DisplayName})
	}
	t.Render()
	return nil
}
