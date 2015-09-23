package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/go-scalingo"
)

func CollaboratorsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	collaborators, err := scalingo.CollaboratorsList(appName)
	if err == nil {

		for _, col := range collaborators {
			fmt.Println(col.Email)
		}
	}

	return nil
}
