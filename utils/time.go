package utils

import (
	"fmt"
	"time"

	"github.com/Scalingo/go-scalingo/v7"
)

const (
	TimeFormat = "2006/01/02 15:04:05"
)

func FormatMaintenanceWindowWithTimezone(maintenanceWindow scalingo.MaintenanceWindow, location *time.Location) string {
	weekdayUTC := time.Weekday(maintenanceWindow.WeekdayUTC)
	weekdayLocal, startingHourLocal := ConvertDayAndHourToTimezone(
		weekdayUTC,
		maintenanceWindow.StartingHourUTC,
		time.UTC,
		location,
	)

	return fmt.Sprintf("%ss at %02d:00 (%02d hours)",
		weekdayLocal.String(),
		startingHourLocal,
		maintenanceWindow.DurationInHour,
	)
}

func ConvertDayAndHourToTimezone(weekday time.Weekday, hour int, inputLocation *time.Location, outputLocation *time.Location) (time.Weekday, int) {
	newTimezoneDate := beginningOfWeek(time.Now().In(inputLocation))

	newTimezoneDate = newTimezoneDate.AddDate(0, 0, int(weekday)-1)
	newTimezoneDate = newTimezoneDate.Add(time.Duration(hour) * time.Hour)

	newTimezoneDate = newTimezoneDate.In(outputLocation)

	return newTimezoneDate.Weekday(), newTimezoneDate.Hour()
}

func beginningOfWeek(t time.Time) time.Time {
	t = beginningOfDay(t)
	weekday := int(t.Weekday())

	weekStartDayInt := int(time.Monday)

	if weekday < weekStartDayInt {
		weekday = weekday + 7 - weekStartDayInt
	} else {
		weekday = weekday - weekStartDayInt
	}
	return t.AddDate(0, 0, -weekday)
}

func beginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
