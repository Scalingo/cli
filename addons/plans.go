package addons

import (
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
)

func Plans(addon string) error {
	plans, err := api.AddonPlansList(addon)
	if err != nil {
		return err
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "Description"})
	for _, plan := range plans {
		t.Append([]string{plan.Name, plan.DisplayName, plan.Description})
	}
	t.Render()
	return nil
}
