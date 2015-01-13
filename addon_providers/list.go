package addon_providers

import (
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List() error {
	addonProviders, err := api.AddonProvidersList()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name"})

	for _, addon := range addonProviders {
		t.Append([]string{addon.NameParam, addon.Name})
	}

	t.Render()
	return nil
}
