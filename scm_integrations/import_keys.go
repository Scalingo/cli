package scm_integrations

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo"
)

func ImportKeys(integration string) error {
	var id string
	var name string
	var keys []scalingo.Key

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	if !utils.IsUUID(integration) {
		i, err := integrationByName(c, integration)
		if err != nil {
			return errgo.Notef(err, "fail to get integration by name")
		}

		id = i.ID
		name = i.ScmType
	} else {
		i, err := integrationByUUID(c, integration)
		if err != nil {
			return errgo.Notef(err, "fail to get integration by uuid")
		}

		id = integration
		name = i.ScmType
	}

	importedKeys, err := c.IntegrationsImportKeys(id)
	if err != nil {
		return errgo.Notef(err, "fail to import keys")
	}

	nbrKeys := len(importedKeys)
	if nbrKeys == 0 {
		alreadyImportedKeys, err := keysContainsName(c, name)
		if err != nil {
			return errgo.Notef(err, "fail to get already imported keys")
		}

		io.Statusf("0 key imported from %s.\n\n", name)
		io.Statusf("You already have %d key(s) that has been previously imported from %s:\n", len(alreadyImportedKeys), name)
		keys = alreadyImportedKeys
	} else {
		keys = importedKeys
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Name", "Content"})
	for _, k := range keys {
		t.Append([]string{k.Name, k.Content[0:20] + "..." + k.Content[len(k.Content)-30:]})
	}
	t.Render()

	if nbrKeys != 0 {
		io.Statusf("%d key(s) have been imported from %s.\n", nbrKeys, name)
	}
	return nil
}

func keysContainsName(c *scalingo.Client, name string) ([]scalingo.Key, error) {
	keys, err := c.KeysList()
	if err != nil {
		return nil, errgo.Notef(err, "fail to get keys")
	}

	var keysAlreadyImported []scalingo.Key

	for _, k := range keys {
		if !strings.Contains(k.Name, name+"-") {
			continue
		}

		if name == "gitlab" && !strings.Contains(k.Name, "gitlab-self-hosted") {
			keysAlreadyImported = append(keysAlreadyImported, k)
		} else if name == "github" && !strings.Contains(k.Name, "github-enterprise") {
			keysAlreadyImported = append(keysAlreadyImported, k)
		}
	}

	return keysAlreadyImported, nil
}
