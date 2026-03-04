package autocomplete

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v10"
)

func AddonsUpgradeAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(c)
	addonName := ""
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	resources, err := client.AddonsList(ctx, appName)
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
			plans, err := client.AddonProviderPlansList(ctx, addonName, scalingo.AddonProviderPlansListOpts{})
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
