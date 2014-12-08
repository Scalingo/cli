package addon_resources

import (
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List(app string) error {
	resources, err := api.AddonResourcesList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Addon", "ID", "Plan"})

	for _, resource := range resources {
		t.Append([]string{resource.Addon, resource.ResourceID, resource.Plan})
	}
	t.Render()

	return nil

}
