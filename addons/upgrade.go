package addons

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Upgrade(ctx context.Context, app, addonID, plan string) error {
	if app == "" {
		return errors.New(ctx, "no app defined")
	} else if addonID == "" {
		return errors.New(ctx, "no addon ID defined")
	} else if plan == "" {
		return errors.New(ctx, "no plan defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	addon, err := checkAddonExists(ctx, c, app, addonID)
	if err != nil {
		return errors.Wrap(ctx, err, "check addon exists")
	}

	planID, err := utils.FindPlan(ctx, c, addon.AddonProvider.ID, plan)
	if err != nil {
		return errors.Wrap(ctx, err, "find plan")
	}

	params, err := c.AddonUpgrade(ctx, app, addon.ID, scalingo.AddonUpgradeParams{
		PlanID: planID,
	})
	if err != nil {
		return errors.Wrap(ctx, err, "addon upgrade")
	}

	io.Status("Addon", addonID, "has been upgraded")
	if len(params.Variables) > 0 {
		io.Info("Modified variables:", params.Variables)
	}
	if len(params.Message) > 0 {
		io.Info("Message from addon provider:", params.Message)
	}
	return nil
}
