package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func Containers(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to list container types")
	}

	containers, err := c.AppsPs(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the application container types")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Amount", "Size", "Command"})

	hasAutoscaler := false
	autoscalers, err := c.AutoscalersList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the autoscalers")
	}

	for _, ct := range containers {
		name := ct.Name

		for _, a := range autoscalers {
			if a.ContainerType == ct.Name {
				hasAutoscaler = true
				name += " (*)"
				break
			}
		}

		amount := fmt.Sprintf("%d", ct.Amount)
		if ct.Command != "" {
			t.Append([]string{name, amount, ct.Size, "`" + ct.Command + "`"})
		} else {
			t.Append([]string{name, amount, ct.Size, "-"})
		}
	}

	t.Render()

	if hasAutoscaler {
		fmt.Println("  (*) has an autoscaler defined")
	}

	return nil
}
