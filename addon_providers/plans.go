package addon_providers

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func Plans(addon string) error {
	c := config.ScalingoUnauthenticatedClient()
	plans, err := c.AddonProviderPlansList(addon)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "Price"})
	for _, plan := range plans {
		t.Append([]string{plan.Name, plan.DisplayName, fmt.Sprintf("%.2f€", plan.Price)})
	}
	t.Render()
	return nil
}
