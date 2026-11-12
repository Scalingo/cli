package dbng

import (
	"context"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func databaseByNameOrID(ctx context.Context, client *scalingo.Client, identifier string) (scalingo.DatabaseNG, error) {
	databases, err := client.Preview().DatabasesList(ctx)
	if err != nil {
		return scalingo.DatabaseNG{}, errors.Wrap(ctx, err, "list databases")
	}

	database, found := findDatabaseByNameOrID(databases, identifier)
	if !found {
		return scalingo.DatabaseNG{}, scalingo.ErrDatabaseNotFound
	}

	return database, nil
}

func findDatabaseByNameOrID(databases []scalingo.DatabaseNG, identifier string) (scalingo.DatabaseNG, bool) {
	for _, database := range databases {
		if database.ID == identifier {
			return database, true
		}
	}

	for _, database := range databases {
		if database.Name == identifier {
			return database, true
		}
	}

	return scalingo.DatabaseNG{}, false
}

func databaseAddon(database scalingo.DatabaseNG, addons []*scalingo.Addon) (*scalingo.Addon, bool) {
	if len(addons) == 1 {
		return addons[0], true
	}

	for _, addon := range addons {
		if addon.ResourceID == database.Database.ResourceID {
			return addon, true
		}
	}

	return nil, false
}

func databaseAndAddonByNameOrID(ctx context.Context, client *scalingo.Client, identifier string) (scalingo.DatabaseNG, *scalingo.Addon, error) {
	database, err := databaseByNameOrID(ctx, client, identifier)
	if err != nil {
		return scalingo.DatabaseNG{}, nil, errors.Wrap(ctx, err, "show database")
	}

	addons, err := client.AddonsList(ctx, database.App.ID)
	if err != nil {
		return scalingo.DatabaseNG{}, nil, errors.Wrap(ctx, err, "list database addons")
	}

	addon, found := databaseAddon(database, addons)
	if !found {
		return scalingo.DatabaseNG{}, nil, errors.New(ctx, "no addon found for database")
	}

	return database, addon, nil
}
