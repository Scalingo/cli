package addons

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

func Upgrade(ctx context.Context, app, addonID, plan string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if addonID == "" {
		return errgo.New("no addon ID defined")
	} else if plan == "" {
		return errgo.New("no plan defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	addon, err := checkAddonExist(ctx, c, app, addonID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	planID, err := utils.FindPlan(ctx, c, addon.AddonProvider.ID, plan)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	params, err := c.AddonUpgrade(ctx, app, addon.ID, scalingo.AddonUpgradeParams{
		PlanID: planID,
	})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
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
