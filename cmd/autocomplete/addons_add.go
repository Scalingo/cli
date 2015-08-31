package autocomplete

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
)

func AddonsAddAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	resources, err := api.AddonsList(appName)
	if len(os.Args) > 1 && err == nil {
		lastArg := os.Args[len(os.Args)-2]
		isAddonNameSet := false

		for _, resource := range resources {
			if lastArg == resource.AddonProvider.ID {
				isAddonNameSet = true
				break
			}
		}

		if isAddonNameSet {
			plans, err := api.AddonProviderPlansList(lastArg)

			if err == nil {
				for _, plan := range plans {
					fmt.Println(plan.Name)
				}
			}
		} else {

			for _, resource := range resources {
				fmt.Println(resource.AddonProvider.ID)
			}
		}
	}

	return nil
}
