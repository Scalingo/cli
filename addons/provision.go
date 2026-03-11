package addons

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Provision(ctx context.Context, app, addon, plan string) error {
	if app == "" {
		return errors.New(ctx, "no app defined")
	} else if addon == "" {
		return errors.New(ctx, "no addon defined")
	} else if plan == "" {
		return errors.New(ctx, "no plan defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	planID, err := utils.FindPlan(ctx, c, addon, plan)
	if err != nil {
		return errors.Wrap(ctx, err, "find plan")
	}

	params, err := c.AddonProvision(ctx, app, scalingo.AddonProvisionParams{
		AddonProviderID: addon,
		PlanID:          planID,
	})
	if err != nil {
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return errors.Wrapf(ctx, err, "Fail to provision addon %v", addon)
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
