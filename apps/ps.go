package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/autoscalers"
	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func Ps(app string) error {
	c := config.ScalingoClient()
	processes, err := c.AppsPs(app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Amount", "Size", "Command"})

	hasAutoscaler := false
	for _, ct := range processes {
		name := ct.Name
		_, err = autoscalers.GetFromContainerType(app, name)
		if err != nil && err != autoscalers.ErrNotFound {
			return errgo.Mask(err, errgo.Any)
		}
		if err == nil {
			hasAutoscaler = true
			name += " (*)"
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
