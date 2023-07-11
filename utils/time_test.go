package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/go-scalingo/v6"
)

func Test_beginningOfDay(t *testing.T) {
	t.Run("returns the beginning of the given date", func(t *testing.T) {
		// Given
		date := time.Date(2023, 3, 15, 22, 50, 12, 13, time.Local)

		// When
		result := beginningOfDay(date)

		// Then
		assert.Equal(t, time.Date(2023, 3, 15, 0, 0, 0, 0, time.Local), result)
	})
}

func Test_BeginningOfWeek(t *testing.T) {
	testDates := []struct {
		description  string
		givenDate    time.Time
		expectedDate time.Time
	}{
		{
			description:  "Wednesday",
			givenDate:    time.Date(2023, 3, 15, 22, 50, 12, 13, time.Local),
			expectedDate: time.Date(2023, 3, 13, 0, 0, 0, 0, time.Local),
		},
		{
			description:  "Monday",
			givenDate:    time.Date(2023, 3, 13, 22, 50, 12, 13, time.Local),
			expectedDate: time.Date(2023, 3, 13, 0, 0, 0, 0, time.Local),
		},
		{
			description:  "Sunday",
			givenDate:    time.Date(2023, 3, 19, 22, 50, 12, 13, time.Local),
			expectedDate: time.Date(2023, 3, 13, 0, 0, 0, 0, time.Local),
		},
	}

	for _, testCase := range testDates {
		t.Run(fmt.Sprintf("returns the beginning of the week for the given date (%s)", testCase.description), func(t *testing.T) {
			// Given
			date := testCase.givenDate

			// When
			result := BeginningOfWeek(date)

			// Then
			assert.Equal(t, testCase.expectedDate, result)
		})
	}
}

func Test_convertUTCDayAndHourToTimezone(t *testing.T) {
	testCases := []struct {
		givenWeekDay        time.Weekday
		givenHour           int
		givenInputLocation  *time.Location
		givenOutputLocation *time.Location
		expectedWeekDay     time.Weekday
		expectedHour        int
		expectedZoneName    string
	}{
		{
			givenWeekDay:        time.Tuesday,
			givenHour:           4,
			givenInputLocation:  time.UTC,
			givenOutputLocation: time.FixedZone("UTC-8", -8*60*60),
			expectedWeekDay:     time.Monday,
			expectedHour:        20,
			expectedZoneName:    "UTC-8",
		},
		{
			givenWeekDay:        time.Tuesday,
			givenHour:           4,
			givenInputLocation:  time.UTC,
			givenOutputLocation: time.FixedZone("UTC+8", 8*60*60),
			expectedWeekDay:     time.Tuesday,
			expectedHour:        12,
			expectedZoneName:    "UTC+8",
		},
		{
			givenWeekDay:        time.Tuesday,
			givenHour:           4,
			givenInputLocation:  time.UTC,
			givenOutputLocation: time.UTC,
			expectedWeekDay:     time.Tuesday,
			expectedHour:        4,
			expectedZoneName:    "UTC",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("converts a weekday / hour from %s to %s", testCase.givenInputLocation.String(), testCase.givenOutputLocation.String()), func(t *testing.T) {
			// Given
			weekDay := testCase.givenWeekDay
			hour := testCase.givenHour
			inputLocation := testCase.givenInputLocation
			outputLocation := testCase.givenOutputLocation

			// When
			resultWeekDay, resultHour := ConvertDayAndHourToTimezone(weekDay, hour, inputLocation, outputLocation)

			// Then
			assert.Equal(t, testCase.expectedWeekDay, resultWeekDay)
			assert.Equal(t, testCase.expectedHour, resultHour)
		})
	}
}

func Test_FormatMaintenanceWindowWithTimezone(t *testing.T) {
	t.Run("it formats the maintenance window in specified timezone", func(t *testing.T) {
		// Given
		maintenanceWindow := scalingo.MaintenanceWindow{
			WeekdayUTC:      1,
			StartingHourUTC: 22,
			DurationInHour:  10,
		}
		location := time.FixedZone("test/zone", 8*60*60)

		// When
		result := FormatMaintenanceWindowWithTimezone(maintenanceWindow, location)

		// Then
		assert.Equal(t, "Tuesdays at 06:00 (10 hours)", result)
	})
}
