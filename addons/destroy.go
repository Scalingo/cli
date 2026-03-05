package addons

import (
	"context"
	"fmt"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Destroy(ctx context.Context, app, addonID string) error {
	if app == "" {
		return errors.New(ctx, "no app defined")
	} else if addonID == "" {
		return errors.New(ctx, "no addon ID defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	addon, err := checkAddonExists(ctx, c, app, addonID)
	if err != nil {
		return errors.Wrap(ctx, err, "check addon exists")
	}

	io.Status("Destroy", addonID)
	io.Warning("This operation is irreversible")
	io.Warning("All related data will be destroyed")
	io.Info("To confirm, type the ID of the addon:")
	fmt.Print("-----> ")

	var validationName string
	fmt.Scan(&validationName)

	if validationName != addonID {
		return errors.Newf(ctx, "'%s' is not '%s', aborting…\n", validationName, addonID)
	}

	err = c.AddonDestroy(ctx, app, addon.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "addon destroy")
	}

	io.Status("Addon", addonID, "has been destroyed")
	return nil
}

func checkAddonExists(ctx context.Context, c *scalingo.Client, app, addonID string) (*scalingo.Addon, error) {
	resources, err := c.AddonsList(ctx, app)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "list addons for app %s", app)
	}
	addonList := []string{}
	for _, r := range resources {
		addonList = append(addonList, r.ID+" ("+r.AddonProvider.Name+")")
		if addonID == r.ID {
			return r, nil
		}
	}
	return nil, errors.Newf(ctx, "Addon "+addonID+" doesn't exist for app "+app+"\nExisting addons:\n  - %v", strings.Join(addonList, "\n  - "))
}
