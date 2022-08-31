package db

import (
	"context"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo/v5"
)

func BackupsConfiguration(ctx context.Context, app, addon string, params scalingo.PeriodicBackupsConfigParams) (scalingo.Database, error) {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return scalingo.Database{}, errgo.Notef(err, "fail to get Scalingo client")
	}
	db, err := client.PeriodicBackupsConfig(ctx, app, addon, params)
	if err != nil {
		return db, errgo.Notef(err, "fail to configure the periodic backups")
	}

	return db, nil
}
