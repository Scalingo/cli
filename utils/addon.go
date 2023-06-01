package utils

import (
	"context"
	"strings"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/go-scalingo/v6"
)

// GetAddonUUIDFromType returns the addon UUID corresponding to the specified application addon type
func GetAddonUUIDFromType(ctx context.Context, addonsClient scalingo.AddonsService, app, addonType string) (string, error) {
	aliases := map[string]string{
		"psql":     "postgresql",
		"pgsql":    "postgresql",
		"postgres": "postgresql",

		"mgo":   "mongodb",
		"mongo": "mongodb",

		"influx": "influxdb",

		"es": "elasticsearch",
	}
	addonTypeAlias, isAlias := aliases[addonType]
	if isAlias {
		addonType = addonTypeAlias
	}

	addons, err := addonsClient.AddonsList(ctx, app)
	if err != nil {
		return "", errors.Notef(ctx, err, "list the addons to get the type UUID")
	}

	for _, addon := range addons {
		if strings.EqualFold(addonType, addon.AddonProvider.Name) {
			return addon.ID, nil
		}
	}

	return "", errors.Newf(ctx, "no '%s' addon exists", addonType)
}
