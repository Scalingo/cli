package integrations

import (
	"fmt"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
	"os"

	"github.com/olekukonko/tablewriter"
)

func ImportKeys(integrationName string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	k, err := integrationByName(c, integrationName)
	if err != nil {
		fmt.Printf("\nfindbyname error\n")
		return errgo.Mask(err)
	}

	fmt.Print("\n\n")

	importedKeys, err := c.IntegrationsImportKeys(k.ID)
	if err != nil {
		fmt.Printf("\nimport keys error\n")
		return errgo.Mask(err)
	}

	fmt.Printf("\n%#v\n", importedKeys)

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Name", "Content"})

	for _, k := range importedKeys {
		t.Append([]string{k.Name, k.Content[0:20] + "..." + k.Content[len(k.Content)-30:]})
	}

	t.Render()

	nbrKeys := len(importedKeys)
	fmt.Printf("%d keys has been imported from %s.\n", nbrKeys, integrationName)
	return nil
}

func integrationByName(c *scalingo.Client, name string) (*scalingo.Integration, error) {
	integrations, err := c.IntegrationsList()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	for _, k := range integrations {
		if k.ScmType == name {
			return &k, nil
		}
	}

	return nil, errgo.New("no such integration '" + name + "'")
}
