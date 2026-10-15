package dbng

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/addons"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Config(ctx context.Context, databaseID string, options addons.UpdateAddonConfigOpts) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, err := client.Preview().DatabaseShow(ctx, databaseID)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	linkedAddons, err := client.AddonsList(ctx, database.App.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "list database addons")
	}

	var addonID string
	for _, addon := range linkedAddons {
		if addon.ResourceID == database.Database.ResourceID {
			addonID = addon.ID
			break
		}
	}
	if addonID == "" && len(linkedAddons) == 1 {
		addonID = linkedAddons[0].ID
	}
	if addonID == "" {
		return errors.New(ctx, "no addon found for database")
	}

	err = addons.UpdateConfig(ctx, database.App.ID, addonID, options)
	if err != nil {
		return errors.Wrap(ctx, err, "update database config")
	}
	fmt.Println("Database config updated.")
	return nil
}
