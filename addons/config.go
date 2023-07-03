package addons

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors/v2"
)

type UpdateAddonConfigOpts struct {
	MaintenanceWindowDay  *string
	MaintenanceWindowHour *int
}

func UpdateConfig(ctx context.Context, app, addon string, options UpdateAddonConfigOpts) error {
	weekdayNameToWeekdayOrderMap := map[string]time.Weekday{
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
		return errors.Notef(ctx, err, "get Scalingo client to update addon config")
	}

	// If addon does not contain a UUID, we consider it contains an addon type (e.g. MongoDB)
	if !strings.HasPrefix(addon, "ad-") {
		addon, err = utils.GetAddonUUIDFromType(ctx, c, app, addon)
		if err != nil {
			return errors.Notef(ctx, err, "fail to get the addon UUID based on its type")
		}
	}

	// fetching the current database maintenance window allows
	// to only overload the specified options
	db, err := c.DatabaseShow(ctx, app, addon)
	if err != nil {
		return errors.Notef(ctx, err, "get database information")
	}

	weekdayLocal, startingHourLocal := utils.ConvertDayAndHourToTimezone(
		time.Weekday(db.MaintenanceWindow.WeekdayUTC),
		db.MaintenanceWindow.StartingHourUTC,
		time.UTC,
		time.Local,
	)

	if options.MaintenanceWindowDay != nil {
		day, ok := weekdayNameToWeekdayOrderMap[strings.ToLower(*options.MaintenanceWindowDay)]
		if !ok {
			return errors.Notef(ctx, err, "invalid weekday '%s'", *options.MaintenanceWindowDay)
		}
		weekdayLocal = day
	}

	if options.MaintenanceWindowHour != nil {
		if *options.MaintenanceWindowHour < 0 || *options.MaintenanceWindowHour > 23 {
			return errors.Notef(ctx, err, "invalid starting hour '%d': it must be between 0 and 23", *options.MaintenanceWindowHour)
		}
		startingHourLocal = *options.MaintenanceWindowHour
	}

	weekdayUTC, startingHourUTC := utils.ConvertDayAndHourToTimezone(weekdayLocal, startingHourLocal, time.Local, time.UTC)

	_, err = c.DatabaseUpdateMaintenanceWindow(ctx, app, addon, scalingo.MaintenanceWindowParams{
		WeekdayUTC:      utils.IntPtr(int(weekdayUTC)),
		StartingHourUTC: utils.IntPtr(startingHourUTC),
	})
	if err != nil {
		return errors.Notef(ctx, err, "update the database maintenance window")
	}

	fmt.Println("Addon config updated.")

	return nil
}
