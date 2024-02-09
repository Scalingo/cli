package users

import (
	"context"
	stdErrors "errors"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

var (
	SupportedAddons                     = []string{"PostgreSQL", "InfluxDB", "MongoDB", "MySQL"}
	ErrDatabaseNotSupportUserManagement = stdErrors.New("Error: DBMS does not support user management")
)

func doesDatabaseHandleUserManagement(ctx context.Context, app, addonUUID string) (bool, error) {
	addonsClient, err := config.ScalingoClient(ctx)
	if err != nil {
		return false, errors.Wrap(ctx, err, "get Scalingo client")
	}

	addon, err := addonsClient.AddonShow(ctx, app, addonUUID)
	if err != nil {
		return false, errors.Wrap(ctx, err, "get the addon to check user management support")
	}

	for _, supportedAddon := range SupportedAddons {
		if strings.EqualFold(supportedAddon, addon.AddonProvider.Name) {
			return true, nil
		}
	}

	return false, nil
}
