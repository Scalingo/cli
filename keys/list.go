package keys

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
	keys, err := c.KeysList()
	if err != nil {
		return errgo.Notef(err, "fail to list SSH keys")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Name", "Content"})

	for _, k := range keys {
		t.Append([]string{k.Name, k.Content[0:20] + "..." + k.Content[len(k.Content)-30:]})
	}

	t.Render()
	return nil
}
