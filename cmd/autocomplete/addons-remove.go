package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
)

func AddonsRemoveAutoComplete(c *cli.Context) error {
	resources, err := api.AddonsList(CurrentAppCompletion(c))
	if err == nil {

		for _, resource := range resources {
			fmt.Println(resource.ResourceID)
		}
	}

	return nil
}
