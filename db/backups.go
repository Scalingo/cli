package db

import (
	"os"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo"
	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	errgo "gopkg.in/errgo.v1"
)

func ListBackups(app, addon string) error {
	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	backups, err := client.BackupList(app, addon)
	if err != nil {
		return errgo.Notef(err, "fail to least backups")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Created At", "Size", "Status"})

	for _, backup := range backups {
		t.Append([]string{
			backup.ID,
			backup.CreatedAt.Format(time.RFC1123),
			humanize.Bytes(backup.Size),
			formatBackupStatus(backup.Status),
		})
	}
	t.Render()
	return nil
}
func formatBackupStatus(status scalingo.BackupStatus) string {
	switch status {
	case scalingo.BackupStatusScheduled:
		return io.Gray(string(status))
	case scalingo.BackupStatusRunning:
		return io.Yellow(string(status))
	case scalingo.BackupStatusDone:
		return io.Green(string(status))
	case scalingo.BackupStatusError:
		return io.BoldRed(string(status))
	}
	return io.BoldRed(string(status))
}
