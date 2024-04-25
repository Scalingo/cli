package db

import (
	"context"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v7"
)

func CreateBackup(ctx context.Context, app, addon string) error {
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Schedule a new backup"
	spinner.Start()
	defer spinner.Stop()

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	backup, err := client.BackupCreate(ctx, app, addon)
	if err != nil {
		return err
	}

	for backup.Status != scalingo.BackupStatusDone &&
		backup.Status != scalingo.BackupStatusError {
		spinner.Lock()
		if backup.Status == scalingo.BackupStatusScheduled {
			spinner.Suffix = " Waiting for the backup to start"
		} else {
			spinner.Suffix = " Waiting for the backup to finish"
		}
		spinner.Unlock()

		time.Sleep(1 * time.Second)

		backup, err = client.BackupShow(ctx, app, addon, backup.ID)
		if err != nil {
			return errgo.Notef(err, "fail to refresh backup state")
		}
	}
	spinner.Stop()

	if backup.Status == scalingo.BackupStatusDone {
		io.Status(color.New(color.FgGreen).Sprint("Backup successfully finished"))
	} else {
		io.Error(color.New(color.FgRed).Sprintf("Backup failed"))
	}
	return nil
}
