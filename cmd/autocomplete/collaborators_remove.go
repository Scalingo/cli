package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func CollaboratorsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	collaborators, err := client.CollaboratorsList(c.Context, appName)
	if err == nil {
		for _, col := range collaborators {
			fmt.Println(col.Email)
		}
	}

	return nil
}
