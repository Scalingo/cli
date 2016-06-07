package addons

import (
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func List(app string) error {
	c := config.ScalingoClient()
	resources, err := c.AddonsList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Addon", "ID", "Plan"})

	for _, resource := range resources {
		t.Append([]string{resource.AddonProvider.Name, resource.ResourceID, resource.Plan.Name})
	}
	t.Render()

	return nil

}
