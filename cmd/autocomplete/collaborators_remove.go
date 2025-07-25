package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func CollaboratorsGenericListAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errors.Wrap(c.Context, err, "fail to get Scalingo client")
	}
	collaborators, err := client.CollaboratorsList(c.Context, appName)
	if err == nil {
		for _, col := range collaborators {
			fmt.Println(col.Email)
		}
	}

	return nil
}
