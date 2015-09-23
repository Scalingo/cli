package autocomplete

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/go-scalingo"
)

func AddonsAddAutoComplete(c *cli.Context) error {
	resources, err := scalingo.AddonProvidersList()
	if len(os.Args) > 1 && err == nil {
		lastArg := os.Args[len(os.Args)-2]
		isAddonNameSet := false

		for _, resource := range resources {
			if lastArg == resource.ID {
				isAddonNameSet = true
				break
			}
		}

		if isAddonNameSet {
			plans, err := scalingo.AddonProviderPlansList(lastArg)

			if err == nil {
				for _, plan := range plans {
					fmt.Println(plan.Name)
				}
			}
		} else {

			for _, resource := range resources {
				fmt.Println(resource.ID)
			}
		}
	}

	return nil
}
