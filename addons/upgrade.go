package addons

import (
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Upgrade(app, resourceID, plan string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if resourceID == "" {
		return errgo.New("no addon ID defined")
	} else if plan == "" {
		return errgo.New("no plan defined")
	}

	addon, err := checkAddonExist(app, resourceID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	planID, err := checkPlanExist(addon.AddonProvider.Name, plan)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	params, err := api.AddonUpgrade(app, addon.ID, planID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Addon", resourceID, "has been upgraded")
	if len(params.Addon.Variables) > 0 {
		io.Info("Modified variables:", params.Addon.Variables)
	}
	if len(params.Addon.Message) > 0 {
		io.Info("Message from addon provider:", params.Addon.Message)
	}
	return nil
}
