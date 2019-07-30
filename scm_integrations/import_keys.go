package scm_integrations

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
)

func ImportKeys(id string) error {
	var keys []scalingo.Key

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integration, err := c.IntegrationsShow(id)
	if err != nil {
		return errgo.Notef(err, "not linked SCM integration or unknown SCM integration")
	}

	importedKeys, err := c.IntegrationsImportKeys(id)
	if err != nil {
		return errgo.Notef(err, "fail to import keys")
	}

	nbrKeys := len(importedKeys)
	if nbrKeys == 0 {
		alreadyImportedKeys, err := keysContainsName(c, integration.ScmType)
		if err != nil {
			return errgo.Notef(err, "fail to list already imported keys")
		}
		alreadyImportedKeysLength := len(alreadyImportedKeys)

		pluralKey := ""
		if alreadyImportedKeysLength > 1 {
			pluralKey = "s"
		}

		io.Statusf("0 key imported from %s.\n\n", integration.ScmType)
		io.Statusf(
			"You already have %d key%s that has been previously imported from %s:\n",
			alreadyImportedKeysLength, pluralKey, integration.ScmType,
		)
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
		pluralKey := ""
		if nbrKeys > 1 {
			pluralKey = "s"
		}
		io.Statusf("%d key%s have been imported from %s.\n", nbrKeys, pluralKey, integration.ScmType)
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

		switch scalingo.SCMType(name) {
		case scalingo.SCMGitlabType:
			if !strings.Contains(k.Name, string(scalingo.SCMGitlabSelfHostedType)) {
				keysAlreadyImported = append(keysAlreadyImported, k)
			}
		case scalingo.SCMGithubType:
			if !strings.Contains(k.Name, string(scalingo.SCMGithubEnterpriseType)) {
				keysAlreadyImported = append(keysAlreadyImported, k)
			}
		case scalingo.SCMGithubEnterpriseType, scalingo.SCMGitlabSelfHostedType:
			keysAlreadyImported = append(keysAlreadyImported, k)
		}
	}

	return keysAlreadyImported, nil
}
