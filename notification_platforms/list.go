package notification_platforms

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List() error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	resources, err := c.NotificationPlatformsList()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name"})

	for _, r := range resources {
		t.Append([]string{r.Name})
	}
	t.Render()

	return nil
}
