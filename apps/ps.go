package apps

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func Ps(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	containers, err := c.AppsContainersPs(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the application containers")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Status", "Command", "Size", "Created At"})

	for _, container := range containers {
		t.Append([]string{container.Label, container.State, container.Command, container.Size, container.CreatedAt.Format(utils.TimeFormat)})
	}
	t.Render()
	return nil
}
