package db

import (
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo/v4"
)

func Show(app, addon string) (scalingo.Database, error) {
	client, err := config.ScalingoClient()
	if err != nil {
		return scalingo.Database{}, errgo.Notef(err, "fail to get Scalingo client")
	}
	db, err := client.DatabaseShow(app, addon)
	if err != nil {
		return db, errgo.Notef(err, "fail to configure the periodic backups")
	}

	return db, nil
}
