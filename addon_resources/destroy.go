package addon_resources

import (
	"fmt"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Destroy(app, resourceID string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if resourceID == "" {
		return errgo.New("no addon ID defined")
	}

	addonResource, err := checkAddonResourceExist(app, resourceID)
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

	err = api.AddonResourceDestroy(app, addonResource.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Addon", resourceID, "has been destroyed")
	return nil
}

func checkAddonResourceExist(app, resourceID string) (*api.AddonResource, error) {
	resources, err := api.AddonResourcesList(app)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	for _, r := range resources {
		if resourceID == r.ResourceID {
			return r, nil
		}
	}
	return nil, errgo.New("Resource ID " + resourceID + " doesn't exist for app " + app)
}
