package addons

import (
	"errors"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
)

func Provision(app, addon, plan string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if addon == "" {
		return errgo.New("no addon defined")
	} else if plan == "" {
		return errgo.New("no plan defined")
	}

	planID, err := checkPlanExist(addon, plan)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	params, err := api.AddonProvision(app, addon, planID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Addon", addon, "has been provisionned")
	io.Info("ID:", params.Addon.ResourceID)
	if len(params.Variables) > 0 {
		io.Info("Modified variables:", params.Variables)
	}
	if len(params.Message) > 0 {
		io.Info("Message from addon provider:", params.Message)
	}
	return nil
}

func checkPlanExist(addon, plan string) (string, error) {
	plans, err := api.AddonProviderPlansList(addon)
	if err != nil {
		return "", errgo.Mask(err, errgo.Any)
	}
	for _, p := range plans {
		if plan == p.Name {
			return p.ID, nil
		}
	}
	return "", errors.New("plan " + plan + " doesn't exist for addon " + addon)
}
