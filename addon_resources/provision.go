package addon_resources

import (
	"errors"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
)

func Provision(app, addon, plan string) error {
	if app == "" {
		return errors.New("no app defined")
	} else if addon == "" {
		return errors.New("no addon defined")
	} else if plan == "" {
		return errors.New("no plan defined")
	}

	planID, err := checkPlanExist(addon, plan)
	if err != nil {
		return err
	}

	params, err := api.AddonResourceProvision(app, addon, planID)
	if err != nil {
		return err
	}

	io.Status("Addon", addon, "has been provisionned")
	io.Info("Resource ID:", params.AddonResource.ResourceID)
	if len(params.Variables) > 0 {
		io.Info("Modified variables:", params.Variables)
	}
	if len(params.Message) > 0 {
		io.Info("Message from addon provider:", params.Message)
	}
	return nil
}

func checkPlanExist(addon, plan string) (string, error) {
	plans, err := api.AddonPlansList(addon)
	if err != nil {
		return "", err
	}
	for _, p := range plans {
		if plan == p.Name {
			return p.ID, nil
		}
	}
	return "", errors.New("plan " + plan + " doesn't exist for addon " + addon)
}
