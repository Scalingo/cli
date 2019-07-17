package db

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	errgo "gopkg.in/errgo.v1"
)

func CreateBackup(app, addon string) error {
	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	err = client.BackupCreate(app, addon)
	if err != nil {
		return err
	}

	io.Status("Successfully ordered to make a backup")
	return nil
}
