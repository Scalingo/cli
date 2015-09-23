package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/go-scalingo"
)

func AddonsPlansAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	resources, err := scalingo.AddonsList(appName)
	if err == nil {
		for _, resource := range resources {
			fmt.Println(resource.AddonProvider.ID)
		}
	}

	return nil
}
