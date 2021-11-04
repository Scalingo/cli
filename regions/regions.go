package regions

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func List() error {
	c, err := config.ScalingoAuthClient()
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo API client")
	}
	regions, err := c.RegionsList()
	if err != nil {
		return errgo.Notef(err, "fail to list available regions")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Name", "Display", "API Endpoint"})

	for _, r := range regions {
		t.Append([]string{r.Name, r.DisplayName, r.API})
	}

	t.Render()
	return nil
}
