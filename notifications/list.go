package notifications

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func List(app string) error {
	c := config.ScalingoClient()
	resources, err := c.NotificationsList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Type", "URL", "Enabled"})

	for _, resource := range resources {
		t.Append([]string{resource.ID, resource.Type, resource.WebHookURL, strconv.FormatBool(resource.Active)})
	}
	t.Render()

	return nil
}
