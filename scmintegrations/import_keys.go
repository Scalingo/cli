package scmintegrations

import (
	"context"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
)

func ImportKeys(ctx context.Context, id string) error {
	var keys []scalingo.Key

	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integration, err := c.SCMIntegrationsShow(ctx, id)
	if err != nil {
		return errgo.Notef(err, "not linked SCM integration or unknown SCM integration")
	}

	importedKeys, err := c.SCMIntegrationsImportKeys(ctx, id)
	if err != nil {
		return errgo.Notef(err, "fail to import keys")
	}

	nbrKeys := len(importedKeys)
	if nbrKeys == 0 {
		alreadyImportedKeys, err := keysContainsName(ctx, c, integration.SCMType.Str())
		if err != nil {
			return errgo.Notef(err, "fail to list already imported keys")
		}
		alreadyImportedKeysLength := len(alreadyImportedKeys)

		pluralKey := ""
		if alreadyImportedKeysLength > 1 {
			pluralKey = "s"
		}

		io.Statusf("0 key imported from %s.\n", scalingo.SCMTypeDisplay[integration.SCMType])
		if alreadyImportedKeysLength == 0 {
			io.Infof("No public key is available in your %s account\n", scalingo.SCMTypeDisplay[integration.SCMType])
			return nil
		}
		io.Info()

		io.Statusf(
			"%d key%s have already been imported from %s:\n",
			alreadyImportedKeysLength, pluralKey, scalingo.SCMTypeDisplay[integration.SCMType],
		)
		keys = alreadyImportedKeys
	} else {
		keys = importedKeys
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Content"})
	for _, k := range keys {
		t.Append([]string{k.Name, k.Content[0:20] + "..." + k.Content[len(k.Content)-30:]})
	}
	t.Render()

	if nbrKeys != 0 {
		pluralKey := ""
		if nbrKeys > 1 {
			pluralKey = "s"
		}
		io.Statusf("%d key%s have been imported from %s.\n", nbrKeys, pluralKey, scalingo.SCMTypeDisplay[integration.SCMType])
	}
	return nil
}

func keysContainsName(ctx context.Context, c *scalingo.Client, name string) ([]scalingo.Key, error) {
	keys, err := c.KeysList(ctx)
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
			if !strings.Contains(k.Name, scalingo.SCMGitlabSelfHostedType.Str()) {
				keysAlreadyImported = append(keysAlreadyImported, k)
			}
		case scalingo.SCMGithubType:
			if !strings.Contains(k.Name, scalingo.SCMGithubEnterpriseType.Str()) {
				keysAlreadyImported = append(keysAlreadyImported, k)
			}
		case scalingo.SCMGithubEnterpriseType, scalingo.SCMGitlabSelfHostedType:
			keysAlreadyImported = append(keysAlreadyImported, k)
		}
	}

	return keysAlreadyImported, nil
}
