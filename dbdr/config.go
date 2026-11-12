package dbdr

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/addons"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Config(ctx context.Context, databaseIdentifier string, options addons.UpdateAddonConfigOpts) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, err := databaseByNameOrID(ctx, client, databaseIdentifier)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	linkedAddons, err := client.AddonsList(ctx, database.App.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "list database addons")
	}

	addon, found := databaseAddon(database, linkedAddons)
	if !found {
		return errors.New(ctx, "no addon found for database")
	}

	err = addons.UpdateConfig(ctx, database.App.ID, addon.ID, options)
	if err != nil {
		return errors.Wrap(ctx, err, "update database config")
	}
	fmt.Println("Database config updated.")
	return nil
}
