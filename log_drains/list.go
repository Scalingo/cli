package log_drains

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

	logDrains, err := c.LogDrainsList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the log drains")
	}

	t := tablewriter.NewWriter(os.Stdout)

	t.SetHeader([]string{"URL"})

	for _, logDrain := range logDrains {
		t.Append([]string{
			logDrain.URL,
		})
	}

	t.Render()
	return nil
}
