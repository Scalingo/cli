package autocomplete

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func AddonsAddAutoComplete(c *cli.Context) error {
	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	resources, err := client.AddonProvidersList()
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
			plans, err := client.AddonProviderPlansList(lastArg)

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
