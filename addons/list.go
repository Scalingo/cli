package addons

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	resources, err := c.AddonsList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Addon", "ID", "Plan", "Status"})

	for _, resource := range resources {
		t.Append([]string{resource.AddonProvider.Name, resource.ID, resource.Plan.Name, string(resource.Status)})
	}
	t.Render()

	return nil

}
