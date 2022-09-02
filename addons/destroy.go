package addons

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v5"
)

func Destroy(ctx context.Context, app, addonID string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if addonID == "" {
		return errgo.New("no addon ID defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	addon, err := checkAddonExist(ctx, c, app, addonID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Destroy", addonID)
	io.Warning("This operation is irreversible")
	io.Warning("All related data will be destroyed")
	io.Info("To confirm, type the ID of the addon:")
	fmt.Print("-----> ")

	var validationName string
	fmt.Scan(&validationName)

	if validationName != addonID {
		return errgo.Newf("'%s' is not '%s', aborting…\n", validationName, addonID)
	}

	err = c.AddonDestroy(ctx, app, addon.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Addon", addonID, "has been destroyed")
	return nil
}

func checkAddonExist(ctx context.Context, c *scalingo.Client, app, addonID string) (*scalingo.Addon, error) {
	resources, err := c.AddonsList(ctx, app)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	addonList := []string{}
	for _, r := range resources {
		addonList = append(addonList, r.ID+" ("+r.AddonProvider.Name+")")
		if addonID == r.ID {
			return r, nil
		}
	}
	return nil, errgo.Newf("Addon "+addonID+" doesn't exist for app "+app+"\nExisting addons:\n  - %v", strings.Join(addonList, "\n  - "))
}
