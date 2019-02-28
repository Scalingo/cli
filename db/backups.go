package db

import (
	"os"
	"strconv"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	errgo "gopkg.in/errgo.v1"
)

func ListBackups(app, addon string) error {
	client := config.ScalingoClient()
	backups, err := client.BackupList(app, addon)
	if err != nil {
		return errgo.Notef(err, "fail to least backups")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "CreatedAt", "Size", "Status"})

	for _, backup := range backups {
		t.Append([]string{
			backup.ID,
			backup.Name,
			backup.CreatedAt.Format(time.RFC1123),
			strconv.FormatUint(backup.Size, 10),
			string(backup.Status),
		})
	}
	t.Render()
	return nil
}
