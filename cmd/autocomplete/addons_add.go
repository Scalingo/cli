package autocomplete

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v9"
)

func AddonsAddAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	resources, err := client.AddonProvidersList(ctx)
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
			plans, err := client.AddonProviderPlansList(ctx, lastArg, scalingo.AddonProviderPlansListOpts{})

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
