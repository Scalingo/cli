package users

import (
	"context"
	"errors"
	"strings"

	"github.com/Scalingo/cli/config"
	scErrors "github.com/Scalingo/go-utils/errors/v2"
)

var ErrDatabaseNotSupportUserManagement = errors.New("Error: DBMS does not support user management")

func doesDatabaseHandleUserManagement(ctx context.Context, app, addonUUID string) (bool, error) {
	supportedAddons := []string{"postgresql", "elasticsearch", "influxdb"}
	addonsClient, err := config.ScalingoClient(ctx)
	if err != nil {
		return false, scErrors.Wrap(ctx, err, "get Scalingo client")
	}

	addon, err := addonsClient.AddonShow(ctx, app, addonUUID)
	if err != nil {
		return false, scErrors.Wrap(ctx, err, "get the addon to check user management support")
	}

	for _, supportedAddon := range supportedAddons {
		if strings.EqualFold(supportedAddon, addon.AddonProvider.Name) {
			return true, nil
		}
	}

	return false, nil
}
