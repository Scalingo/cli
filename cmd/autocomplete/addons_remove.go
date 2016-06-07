package autocomplete

import (
	"fmt"

	"github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/config"
)

func AddonsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client := config.ScalingoClient()
	resources, err := client.AddonsList(appName)
	if err == nil {

		for _, resource := range resources {
			fmt.Println(resource.ResourceID)
		}
	}

	return nil
}
