package db

import (
	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo"
	errgo "gopkg.in/errgo.v1"
)

func BackupsConfiguration(app, addon string, params scalingo.PeriodicBackupsConfigParams) (scalingo.Database, error) {
	client, err := config.ScalingoClient()
	if err != nil {
		return scalingo.Database{}, errgo.Notef(err, "fail to get Scalingo client")
	}
	db, err := client.PeriodicBackupsConfig(app, addon, params)
	if err != nil {
		return db, errgo.Notef(err, "fail to configure the periodic backups")
	}

	return db, nil
}
