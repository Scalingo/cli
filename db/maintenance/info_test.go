package maintenance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v6"
)

func TestGetNextMaintenanceLocalWindow(t *testing.T) {
	functionExecutionTime := dateOfTheWeekUTC(time.Wednesday, 8)

	t.Run("when the maintenance is notified", func(t *testing.T) {
		maintenanceType := scalingo.MaintenanceStatusNotified

		t.Run("when the maintenance should be executed next week", func(t *testing.T) {
			// given
			maintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Sunday, 0)
			expectedMaintenanceExecutionTime := maintenanceWindowExecutionTime.AddDate(0, 0, 7)
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(maintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: maintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})

		t.Run("when the maintenance window is ongoing", func(t *testing.T) {
			// given
			expectedMaintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Wednesday, 7)
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(expectedMaintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: expectedMaintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceWindowExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceWindowExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})

		t.Run("when the maintenance window is this week", func(t *testing.T) {
			// given
			expectedMaintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Saturday, 7)
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(expectedMaintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: expectedMaintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceWindowExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceWindowExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})
	})

	t.Run("when the maintenance is queued", func(t *testing.T) {
		maintenanceType := scalingo.MaintenanceStatusQueued

		t.Run("when the maintenance should be executed next week", func(t *testing.T) {
			// given
			maintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Sunday, 0)
			expectedMaintenanceExecutionTime := maintenanceWindowExecutionTime.AddDate(0, 0, 7)
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(maintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: maintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})

		t.Run("when the maintenance window is ongoing", func(t *testing.T) {
			// given
			expectedMaintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Wednesday, 7)
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(expectedMaintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: expectedMaintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceWindowExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceWindowExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})

		t.Run("when the maintenance window is this week", func(t *testing.T) {
			// given
			maintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Saturday, 7)
			expectedMaintenanceExecutionTime := maintenanceWindowExecutionTime
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(maintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: maintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})
	})

	t.Run("when the maintenance is scheduled", func(t *testing.T) {
		maintenanceType := scalingo.MaintenanceStatusScheduled

		t.Run("when the maintenance should be executed next week", func(t *testing.T) {
			// given
			maintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Sunday, 0)
			expectedMaintenanceExecutionTime := maintenanceWindowExecutionTime.AddDate(0, 0, 7)
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(maintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: maintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})

		t.Run("when the maintenance window is ongoing it should be executed next week", func(t *testing.T) {
			// given
			maintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Wednesday, 7)
			expectedMaintenanceExecutionTime := maintenanceWindowExecutionTime.AddDate(0, 0, 7)
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(maintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: maintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})

		t.Run("when the maintenance window is this week", func(t *testing.T) {
			// given
			maintenanceWindowExecutionTime := dateOfTheWeekUTC(time.Saturday, 7)
			expectedMaintenanceExecutionTime := maintenanceWindowExecutionTime
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(maintenanceWindowExecutionTime.Weekday()),
				StartingHourUTC: maintenanceWindowExecutionTime.Hour(),
				DurationInHour:  8,
			}

			// when
			nextExec, untilExec := getNextLocalMaintenanceWindow(functionExecutionTime, maintenanceWindow, maintenanceType)

			// Then
			assert.Equal(t, expectedMaintenanceExecutionTime.Local(), nextExec)
			assert.Equal(t, expectedMaintenanceExecutionTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
		})
	})
}

func dateOfTheWeekUTC(weekday time.Weekday, hour int) time.Time {
	date := utils.BeginningOfWeek(time.Now().In(time.UTC))
	date = date.AddDate(0, 0, int(weekday)-1)
	date = date.Add(time.Duration(hour) * time.Hour)

	return date
}
