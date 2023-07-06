package utils

import (
	"context"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func GetAddonUUIDFromType(ctx context.Context, app, addonTypeOrUUID string) (string, error) {
	// If addon does not contain a UUID, we consider it contains an addon type (e.g. MongoDB)
	if strings.HasPrefix(addonTypeOrUUID, "ad-") {
		return addonTypeOrUUID, nil
	}
	addonType := addonTypeOrUUID

	addonsClient, err := config.ScalingoClient(ctx)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "unable to get Scalingo client")
	}

	aliases := map[string]string{
		"psql":     "postgresql",
		"pgsql":    "postgresql",
		"postgres": "postgresql",

		"mgo":   "mongodb",
		"mongo": "mongodb",

		"influx": "influxdb",

		"es": "elasticsearch",
	}
	addonTypeFromAlias, isAlias := aliases[addonType]
	if isAlias {
		addonType = addonTypeFromAlias
	}

	addons, err := addonsClient.AddonsList(ctx, app)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "list the addons to get the type UUID")
	}

	for _, addon := range addons {
		if strings.EqualFold(addonType, addon.AddonProvider.Name) {
			return addon.ID, nil
		}
	}

	return "", errors.Newf(ctx, "no '%s' addon exists", addonType)
}
