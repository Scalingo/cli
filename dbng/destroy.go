package dbng

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Destroy(ctx context.Context, databaseIdentifier string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	databases, err := c.Preview().DatabasesList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list databases")
	}

	db, found := findDatabaseByNameOrID(databases, databaseIdentifier)
	if !found {
		return errors.Wrap(ctx, scalingo.ErrDatabaseNotFound, "show database")
	}

	io.Warningf("You're going to delete database %s ('%s'),\n", db.ID, db.Name)
	io.Warning()
	io.Warning("This operation is irreversible, all data including backups of your database will be deleted.")

	fmt.Print("\nTo confirm type the ID or the name of the database: ")
	var validationID string
	_, err = fmt.Scan(&validationID)
	if err != nil {
		return errors.Wrap(ctx, err, "delete database confirmation")
	}
	fmt.Println()

	if validationID != db.ID && validationID != db.Name {
		return errors.Newf(ctx, "'%s' is not the ID or the name of the database, aborting…\n", validationID)
	}

	err = c.Preview().DatabaseDestroy(ctx, db.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy database")
	}

	io.Statusf("The database %s has been deleted.\n", db.ID)

	return nil
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
