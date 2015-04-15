package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
)

func Ps(app string) error {
	processes, err := api.AppsPs(app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Amount", "Size", "Command"})

	for _, ct := range processes {
		amount := fmt.Sprintf("%d", ct.Amount)
		if ct.Command != "" {
			t.Append([]string{ct.Name, amount, ct.Size, "`" + ct.Command + "`"})
		} else {
			t.Append([]string{ct.Name, amount, ct.Size, "-"})
		}
	}

	t.Render()
	return nil
}
