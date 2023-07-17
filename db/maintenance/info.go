package maintenance

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Info(ctx context.Context, app, addonUUID, maintenanceID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Notef(ctx, err, "get Scalingo client")
	}

	maintenanceInfo, err := c.DatabaseShowMaintenance(ctx, app, addonUUID, maintenanceID)
	if err != nil {
		return errors.Notef(ctx, err, "get the maintenance")
	}

	var startedAtMessage string
	var endedAt string

	if maintenanceInfo.StartedAt != nil {
		startedAtMessage = maintenanceInfo.StartedAt.Local().Format(utils.TimeFormat)
	} else {
		database, err := c.DatabaseShow(ctx, app, addonUUID)
		if err != nil {
			return errors.Notef(ctx, err, "get database information")
		}

		nextMaintenanceWindowStartingDate, nextMaintenanceWindowEndingDate := getNextLocalMaintenanceWindow(time.Now(), database.MaintenanceWindow, maintenanceInfo.Status)

		now := time.Now()
		startedAtMessage = fmt.Sprintf("Next Maintenance Window: %s", nextMaintenanceWindowStartingDate.Format(utils.TimeFormat))
		if nextMaintenanceWindowStartingDate.Before(now) && nextMaintenanceWindowEndingDate.After(now) && maintenanceInfo.Status == scalingo.MaintenanceStatusNotified {
			startedAtMessage = fmt.Sprintf("Could be executed during current maintenance window:\n %s - %s",
				nextMaintenanceWindowStartingDate.Format(utils.TimeFormat),
				nextMaintenanceWindowEndingDate.Format("15:04:00"))
		}
	}

	if maintenanceInfo.EndedAt != nil {
		endedAt = maintenanceInfo.EndedAt.Local().Format(utils.TimeFormat)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetAutoWrapText(false)
	t.Append([]string{"ID", fmt.Sprintf("%v", maintenanceInfo.ID)})
	t.Append([]string{"Type", fmt.Sprintf("%v", maintenanceInfo.Type)})
	t.Append([]string{"Started At", fmt.Sprintf("%v", startedAtMessage)})
	t.Append([]string{"Ended At", fmt.Sprintf("%v", endedAt)})
	t.Append([]string{"Status", fmt.Sprintf("%v", maintenanceInfo.Status)})

	t.Render()

	return nil
}

// getNextLocalMaintenanceWindow returns the boundaries of the next maintenance date.
// The logic is as follow:
// We make the difference between the Time.Now and the date of the maintenance in the week.
// Is the next maintenance windows is before now:
//
//	Yes => then the next maintenance is the day of the maintenance + 1 week
//			except in the case if we are already in the range of the maintenance window and status different than "Scheduled".
//	No 	=> then the next maintenance is the maintenance time from the database information.
//
// It returns the current maintenance if we are within the maintenance boundaries or returns the next maintenance.
func getNextLocalMaintenanceWindow(now time.Time, maintenanceWindow scalingo.MaintenanceWindow, status scalingo.MaintenanceStatus) (time.Time, time.Time) {
	var nextMaintenanceWindowStartingDate time.Time
	now = now.UTC()

	dayDiff := int(time.Weekday(maintenanceWindow.WeekdayUTC) - now.Weekday())
	year, month, date := now.AddDate(0, 0, dayDiff).Date()
	nextMaintenanceWindowStartingDate = time.Date(year, month, date, maintenanceWindow.StartingHourUTC, 0, 0, 0, now.Location())

	nextMaintenanceWindowEndingDate := nextMaintenanceWindowStartingDate.Add(time.Hour * time.Duration(maintenanceWindow.DurationInHour))

	// if we are in a maintenance window and if the maintenance status is `Scheduled` then the maintenance will be executed next time
	if nextMaintenanceWindowStartingDate.Before(now) && (nextMaintenanceWindowEndingDate.Before(now) || status == scalingo.MaintenanceStatusScheduled) {
		nextMaintenanceWindowStartingDate = nextMaintenanceWindowStartingDate.AddDate(0, 0, 7)
	}

	return nextMaintenanceWindowStartingDate.Local(), nextMaintenanceWindowStartingDate.Add(time.Hour * time.Duration(maintenanceWindow.DurationInHour)).Local()
}
