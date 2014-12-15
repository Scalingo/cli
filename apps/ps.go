package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func Ps(app string) error {
	processes, err := api.AppsPs(app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Amount", "Command"})

	for _, ct := range processes {
		amount := fmt.Sprintf("%d", ct.Amount)
		if ct.Command != "" {
			t.Append([]string{ct.Name, amount, "`" + ct.Command + "`"})
		} else {
			t.Append([]string{ct.Name, amount, "-"})
		}
	}

	t.Render()
	return nil
}
