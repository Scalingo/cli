package integrations

import (
	"fmt"
	"gopkg.in/errgo.v1"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/olekukonko/tablewriter"
)

func ImportKeys(integration string) error {
	var id string
	var name string

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	if !utils.IsUUID(integration) {
		i, err := integrationByName(c, integration)
		if err != nil {
			return errgo.Mask(err)
		}

		id = i.ID
		name = i.ScmType
	} else {
		i, err := integrationByUUID(c, integration)
		if err != nil {
			return errgo.Mask(err)
		}

		id = integration
		name = i.ScmType
	}

	importedKeys, err := c.IntegrationsImportKeys(id)
	if err != nil {
		return errgo.Mask(err)
	}

	nbrKeys := len(importedKeys)
	if nbrKeys == 0 {
		fmt.Printf("0 keys imported from %s.\n", name)
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Name", "Content"})
	for _, k := range importedKeys {
		t.Append([]string{k.Name, k.Content[0:20] + "..." + k.Content[len(k.Content)-30:]})
	}
	t.Render()

	fmt.Printf("%d keys has been imported from %s.\n", nbrKeys, name)
	return nil
}
