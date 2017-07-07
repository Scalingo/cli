package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/Scalingo/cli/config"
)

func CollaboratorsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client := config.ScalingoClient()
	collaborators, err := client.CollaboratorsList(appName)
	if err == nil {

		for _, col := range collaborators {
			fmt.Println(col.Email)
		}
	}

	return nil
}
