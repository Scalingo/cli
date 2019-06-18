package integrations

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func List() error {
	/*
	c := config.ScalingoClient()
	keys, err := c.KeysList()
	if err != nil {
		return errgo.Mask(err)
	}
	*/

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Type", "URL", "Username", "Email"})

	/*
	for _, k := range keys {
		t.Append([]string{k.Name, k.Content[0:20] + "..." + k.Content[len(k.Content)-30:]})
	}*/

	t.Append([]string{"GitLab", "", "brandon-welsch", "dev@brandon-welsch.eu"})
	t.Append([]string{"GitHub Enterprise", "https://ghe.scalingo.com", "brandon", "dev@brandon-welsch.eu"})
	t.Append([]string{"GitLab Self Hosted", "https://gitlab.sysroot.ovh", "brandon_welsch", "dev@brandon-welsch.eu"})

	t.Render()
	return nil
}
