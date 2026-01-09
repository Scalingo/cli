package utils

import (
	"context"

	stderrors "github.com/pkg/errors"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v9/debug"
	"github.com/Scalingo/go-utils/errors/v2"
)

var ErrResourceNotFound = stderrors.New("resource name not found")

const databasesResourceFlag = "dedicated-database"

// IsResourceDatabase returns true if the current resource is a database. It returns
// an ErrResourceNotFound error if no database name has been found.
func IsResourceDatabase(ctx context.Context, resourceName string) (bool, error) {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return false, errors.Wrap(ctx, err, "get Scalingo client")
	}

	app, err := client.AppsShow(ctx, resourceName)
	if err != nil {
		debug.Println("[detect] Unable to show app to check if resource is a database:", err)
		return false, nil
	}

	for f, v := range app.Flags {
		if f == databasesResourceFlag && v {
			return true, nil
		}
	}

	return false, ErrResourceNotFound
}
