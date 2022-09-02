package addon_providers

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

	addonProviders, err := c.AddonProvidersList(ctx)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name"})

	for _, addon := range addonProviders {
		t.Append([]string{addon.ID, addon.Name})
	}

	t.Render()
	return nil
}
