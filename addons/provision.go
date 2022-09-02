package addons

import (
	"context"
	"errors"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v5"
)

func Provision(ctx context.Context, app, addon, plan string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if addon == "" {
		return errgo.New("no addon defined")
	} else if plan == "" {
		return errgo.New("no plan defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	planID, err := checkPlanExist(ctx, c, addon, plan)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	params, err := c.AddonProvision(ctx, app, scalingo.AddonProvisionParams{
		AddonProviderID: addon,
		PlanID:          planID,
	})
	if err != nil {
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return errgo.Notef(err, "Fail to provision addon %v", addon)
		}
		// If error is Payment Required and user tries to exceed its free trial
		return utils.AskAndStopFreeTrial(ctx, c, func() error {
			return Provision(ctx, app, addon, plan)
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

func checkPlanExist(ctx context.Context, c *scalingo.Client, addon, plan string) (string, error) {
	plans, err := c.AddonProviderPlansList(ctx, addon)
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
