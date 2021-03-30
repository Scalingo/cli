package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func ContainerTypes(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to list container types")
	}

	containerTypes, err := c.AppsContainerTypes(app)
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

	for _, containerType := range containerTypes {
		name := containerType.Name

		for _, a := range autoscalers {
			if a.ContainerType == containerType.Name {
				hasAutoscaler = true
				name += " (*)"
				break
			}
		}

		amount := fmt.Sprintf("%d", containerType.Amount)
		if containerType.Command != "" {
			t.Append([]string{name, amount, containerType.Size, "`" + containerType.Command + "`"})
		} else {
			t.Append([]string{name, amount, containerType.Size, "-"})
		}
	}

	t.Render()

	if hasAutoscaler {
		fmt.Println("  (*) has an autoscaler defined")
	}

	return nil
}
