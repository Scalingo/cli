package maintenance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-scalingo/v9"
)

type test struct {
	maintenanceStatuses []scalingo.MaintenanceStatus
	expectedRunTime     time.Time
	maintenanceRunDate  time.Time
}

func TestGetNextMaintenanceLocalWindow(t *testing.T) {
	now := dateAt(t, "Wednesday, 19-Jul-23 08:00:00 UTC")

	tests := map[string]test{
		"when the maintenance is notified, queued or scheduled, the maintenance should be executed a week later when execution date is in the past": {
			maintenanceStatuses: []scalingo.MaintenanceStatus{scalingo.MaintenanceStatusNotified, scalingo.MaintenanceStatusQueued, scalingo.MaintenanceStatusScheduled},
			maintenanceRunDate:  dateAt(t, "Sunday, 16-Jul-23 00:00:00 UTC"),
			expectedRunTime:     dateAt(t, "Sunday, 23-Jul-23 00:00:00 UTC"),
		},
		"when the maintenance is notified or queued, the maintenance should be ongoing when execution date is within maintenance boundaries": {
			maintenanceStatuses: []scalingo.MaintenanceStatus{scalingo.MaintenanceStatusNotified, scalingo.MaintenanceStatusQueued},
			maintenanceRunDate:  dateAt(t, "Wednesday, 19-Jul-23 07:00:00 UTC"),
			expectedRunTime:     dateAt(t, "Wednesday, 19-Jul-23 07:00:00 UTC"),
		},
		"when the maintenance is scheduled, the maintenance should be executed when execution date is within maintenance boundaries": {
			maintenanceStatuses: []scalingo.MaintenanceStatus{scalingo.MaintenanceStatusScheduled},
			maintenanceRunDate:  dateAt(t, "Wednesday, 19-Jul-23 07:00:00 UTC"),
			expectedRunTime:     dateAt(t, "Saturday, 26-Jul-23 07:00:00 UTC"),
		},
		"when the maintenance is notified, queued or scheduled, the maintenance should be executed this week when execution date in future": {
			maintenanceStatuses: []scalingo.MaintenanceStatus{scalingo.MaintenanceStatusNotified, scalingo.MaintenanceStatusQueued, scalingo.MaintenanceStatusScheduled},
			maintenanceRunDate:  dateAt(t, "Saturday, 22-Jul-23 07:00:00 UTC"),
			expectedRunTime:     dateAt(t, "Saturday, 22-Jul-23 07:00:00 UTC"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			maintenanceWindow := scalingo.MaintenanceWindow{
				WeekdayUTC:      int(test.maintenanceRunDate.Weekday()),
				StartingHourUTC: test.maintenanceRunDate.Hour(),
				DurationInHour:  8,
			}

			for _, status := range test.maintenanceStatuses {
				// when
				nextExec, untilExec := getNextLocalMaintenanceWindow(now, maintenanceWindow, status)

				// Then
				assert.Equal(t, test.expectedRunTime.Local(), nextExec)
				assert.Equal(t, test.expectedRunTime.Local().Add(time.Duration(maintenanceWindow.DurationInHour)*time.Hour), untilExec)
			}
		})
	}
}

func dateAt(t *testing.T, date string) time.Time {
	parsedTime, err := time.Parse(time.RFC850, date)

	require.NoError(t, err)

	return parsedTime
}
