package addons

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
)

type ConfigOpts struct {
	MaintenanceWindowDay  *string
	MaintenanceWindowHour *int
}

func Config(ctx context.Context, app, addon string, options ConfigOpts) error {
	weekdaysMap := map[string]time.Weekday{
		strings.ToLower(time.Sunday.String()):    time.Sunday,
		strings.ToLower(time.Monday.String()):    time.Monday,
		strings.ToLower(time.Tuesday.String()):   time.Tuesday,
		strings.ToLower(time.Wednesday.String()): time.Wednesday,
		strings.ToLower(time.Thursday.String()):  time.Thursday,
		strings.ToLower(time.Friday.String()):    time.Friday,
		strings.ToLower(time.Saturday.String()):  time.Saturday,
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "get Scalingo client to update addon config")
	}

	// fetching the current database maintenance window allows
	// to only overload the specified options
	db, err := c.DatabaseShow(ctx, app, addon)
	if err != nil {
		return errgo.Notef(err, "fail to get database information")
	}

	weekdayLocal, startingHourLocal := utils.ConvertDayAndHourToTimezone(
		time.Weekday(db.MaintenanceWindow.WeekdayUTC),
		db.MaintenanceWindow.StartingHourUTC,
		time.UTC,
		time.Local,
	)

	if options.MaintenanceWindowDay != nil {
		day, ok := weekdaysMap[strings.ToLower(*options.MaintenanceWindowDay)]
		if !ok {
			return errgo.Notef(err, "invalid weekday '%s'", *options.MaintenanceWindowDay)
		}
		weekdayLocal = day
	}

	if options.MaintenanceWindowHour != nil {
		if *options.MaintenanceWindowHour < 0 || *options.MaintenanceWindowHour > 23 {
			return errgo.Notef(err, "invalid starting hour '%d': it should be between 0 and 23", *options.MaintenanceWindowHour)
		}
		startingHourLocal = *options.MaintenanceWindowHour
	}

	weekdayUTC, startingHourUTC := utils.ConvertDayAndHourToTimezone(weekdayLocal, startingHourLocal, time.Local, time.UTC)

	_, err = c.DatabaseUpdateMaintenanceWindow(ctx, app, addon, scalingo.MaintenanceWindowParams{
		WeekdayUTC:      utils.IntPtr(int(weekdayUTC)),
		StartingHourUTC: utils.IntPtr(startingHourUTC),
	})
	if err != nil {
		return errgo.Notef(err, "fail to update the application information")
	}

	fmt.Println("Addon config updated.")

	return nil
}
