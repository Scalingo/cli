package addon_resources

import (
	"errors"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
)

func Upgrade(app, resourceID, plan string) error {
	if app == "" {
		return errors.New("no app defined")
	} else if resourceID == "" {
		return errors.New("no addon ID defined")
	} else if plan == "" {
		return errors.New("no plan defined")
	}

	addonResource, err := checkAddonResourceExist(app, resourceID)
	if err != nil {
		return err
	}

	planID, err := checkPlanExist(addonResource.Addon, plan)
	if err != nil {
		return err
	}

	params, err := api.AddonResourceUpgrade(app, addonResource.ID, planID)
	if err != nil {
		return err
	}

	io.Status("Addon", resourceID, "has been upgraded")
	if len(params.Variables) > 0 {
		io.Info("Modified variables:", params.Variables)
	}
	if len(params.Message) > 0 {
		io.Info("Message from addon provider:", params.Message)
	}
	return nil
}
