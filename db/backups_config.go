package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo/v5"
)

// BackupsConfiguration is the handler which handles the update of backup
// configuration for a database addon
func BackupsConfiguration(ctx context.Context, app, addon string, params scalingo.PeriodicBackupsConfigParams) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	db, err := client.DatabaseShow(ctx, app, addon)
	if err != nil {
		return errgo.Notef(err, "fail to get database info")
	}

	if params.ScheduledAt != nil && len(db.PeriodicBackupsScheduledAt) > 1 {
		io.Warning("Your database is backed up multiple times a day at " +
			formatScheduledAt(db.PeriodicBackupsScheduledAt) +
			". Please ask the support to update the frequency of these backups.")
		return nil
	}

	db, err = client.PeriodicBackupsConfig(ctx, app, addon, params)
	if err != nil {
		return errgo.Notef(err, "fail to configure the periodic backups")
	}

	if db.PeriodicBackupsEnabled {
		io.Statusf("Periodic backups will be done daily at %s\n", formatScheduledAt(db.PeriodicBackupsScheduledAt))
	} else {
		io.Status("Periodic backups are disabled")
	}

	return nil
}

func formatScheduledAt(hours []int) string {
	hoursStr := make([]string, len(hours))
	for i, h := range hours {
		hUTC := time.Date(1986, 7, 22, h, 0, 0, 0, time.UTC)
		hLocal := hUTC.In(time.Local)
		hoursStr[i] = strconv.Itoa(hLocal.Hour())
	}

	tz, _ := time.Now().In(time.Local).Zone()
	return fmt.Sprintf("%s:00 %s", strings.Join(hoursStr, ":00, "), tz)
}
