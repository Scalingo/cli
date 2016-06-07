package autocomplete

import (
	"fmt"
	"os"

	"github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/config"
)

func AddonsUpgradeAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	addonName := ""
	if appName == "" {
		return nil
	}

	client := config.ScalingoClient()
	resources, err := client.AddonsList(appName)
	if len(os.Args) > 1 && err == nil {

		lastArg := os.Args[len(os.Args)-2]
		isAddonIDSet := false
		for _, resource := range resources {
			if lastArg == resource.ResourceID {
				isAddonIDSet = true
				addonName = resource.AddonProvider.ID
				break
			}
		}

		if isAddonIDSet && addonName != "" {
			plans, err := client.AddonProviderPlansList(addonName)
			if err == nil {
				for _, plan := range plans {
					fmt.Println(plan.Name)
				}
			}
		} else {
			for _, resource := range resources {
				fmt.Println(resource.ResourceID)
			}
		}
	}

	return nil
}
