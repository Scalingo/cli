package autocomplete

import (
	"context"
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-utils/errors/v3"
)

func AddonsAddAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
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
