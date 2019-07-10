package integrations

import (
	"fmt"
	"gopkg.in/errgo.v1"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
)

func List() error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integrations, err := c.IntegrationsList()
	if err != nil {
		return errgo.Mask(err)
	}

	nbrKeys := len(integrations)
	if nbrKeys == 0 {
		fmt.Printf("0 integrations linked to your scalingo account\n")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"ID", "Type", "URL", "Username", "Email"})
	for _, i := range integrations {
		t.Append([]string{i.ID, i.ScmType, i.Url, i.Username, i.Email})
	}
	t.Render()
	return nil
}
