package addon_providers

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func List() error {
	c := config.ScalingoUnauthenticatedClient()
	addonProviders, err := c.AddonProvidersList()
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
