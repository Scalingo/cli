package addon_providers

import (
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/go-scalingo"
)

func Plans(addon string) error {
	plans, err := scalingo.AddonProviderPlansList(addon)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "Description"})
	for _, plan := range plans {
		t.Append([]string{plan.Name, plan.DisplayName, plan.ShortDescription})
	}
	t.Render()
	return nil
}
