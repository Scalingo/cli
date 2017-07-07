package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/Scalingo/cli/config"
)

func AddonsPlansAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client := config.ScalingoClient()
	resources, err := client.AddonsList(appName)
	if err == nil {
		for _, resource := range resources {
			fmt.Println(resource.AddonProvider.ID)
		}
	}

	return nil
}
