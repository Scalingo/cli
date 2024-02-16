package addonproviders

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

	addonProviders, err := c.AddonProvidersList(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "list addon providers")
	}
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name"})

	for _, addon := range addonProviders {
		t.Append([]string{addon.ID, addon.Name})
	}

	t.Render()
	return nil
}
