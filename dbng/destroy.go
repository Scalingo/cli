package dbng

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Destroy(ctx context.Context, appID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	db, err := c.Preview().DatabaseShow(ctx, appID)
	if err != nil {
		return errors.Wrap(ctx, err, "delete database")
	}

	io.Warningf("You're going to delete database %s ('%s'),\n", db.DatabaseInfo.ID, db.DatabaseInfo.Name)
	io.Warning()
	io.Warning("This operation is irreversible, all data including backups of your database will be deleted.")

	fmt.Print("\nTo confirm type the ID or the name of the database: ")
	var validationID string
	_, err = fmt.Scan(&validationID)
	if err != nil {
		return errors.Wrap(ctx, err, "delete database confirmation")
	}
	fmt.Println()

	if validationID != db.DatabaseInfo.ID && validationID != db.DatabaseInfo.Name {
		return errors.Newf(ctx, "'%s' is not the ID or the name of the database, aborting…\n", validationID)
	}

	err = c.Preview().DatabaseDestroy(ctx, appID)
	if err != nil {
		return errors.Wrap(ctx, err, "delete database")
	}

	io.Statusf("The database %s has been deleted.", db.DatabaseInfo.ID)

	return nil
}
