package addons

import (
	"errors"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"gopkg.in/errgo.v1"
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

	c := config.ScalingoClient()
	params, err := c.AddonProvision(app, addon, planID)
	if err != nil {
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return errgo.Notef(err, "Fail to provision addon %v", addon)
		}
		// If error is Payment Required and user tries to exceed its free trial
		return utils.AskAndStopFreeTrial(c, func() error {
			return Provision(app, addon, plan)
		})
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
	c := config.ScalingoClient()
	plans, err := c.AddonProviderPlansList(addon)
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
