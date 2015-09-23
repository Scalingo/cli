package addons

import (
	"fmt"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
)

func Destroy(app, resourceID string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if resourceID == "" {
		return errgo.New("no addon ID defined")
	}

	addon, err := checkAddonExist(app, resourceID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Destroy", resourceID)
	io.Warning("This operation is irreversible")
	io.Warning("All related data will be destroyed")
	io.Info("To confirm, type the ID of the addon:")
	fmt.Print("-----> ")

	var validationName string
	fmt.Scan(&validationName)

	if validationName != resourceID {
		return errgo.Newf("'%s' is not '%s', aborting…\n", validationName, resourceID)
	}

	err = scalingo.AddonDestroy(app, addon.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Addon", resourceID, "has been destroyed")
	return nil
}

func checkAddonExist(app, resourceID string) (*scalingo.Addon, error) {
	resources, err := scalingo.AddonsList(app)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	addonList := []string{}
	for _, r := range resources {
		addonList = append(addonList, r.ResourceID+" ("+r.AddonProvider.Name+")")
		if resourceID == r.ResourceID {
			return r, nil
		}
	}
	return nil, errgo.Newf("Addon "+resourceID+" doesn't exist for app "+app+"\nExisting addons:\n  - %v", strings.Join(addonList, "\n  - "))
}
