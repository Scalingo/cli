package addon_providers

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Plans(addon string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	plans, err := c.AddonProviderPlansList(addon)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "Price/month"})
	for _, plan := range plans {
		t.Append([]string{plan.Name, plan.DisplayName, fmt.Sprintf("%.2f€", plan.Price)})
	}
	t.Render()
	return nil
}
