package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseScheduleAtFlag(t *testing.T) {
	tests := map[string]struct {
		flag          string
		expectedHour  int
		expectedLoc   string
		expectedError string
	}{
		"single hour uses local timezone": {
			flag:         "3",
			expectedHour: 3,
			expectedLoc:  time.Local.String(),
		},
		"hour and timezone": {
			flag:         "3 Europe/Paris",
			expectedHour: 3,
			expectedLoc:  "Europe/Paris",
		},
		"hour with minutes and timezone (1)": {
			flag:         "4:00 Europe/Paris",
			expectedHour: 4,
			expectedLoc:  "Europe/Paris",
		},
		"hour with minutes and timezone (2)": {
			flag:         "4h00 Europe/Paris",
			expectedHour: 4,
			expectedLoc:  "Europe/Paris",
		},
		"hour with minutes and timezone (3)": {
			flag:         "4H00 Europe/Paris",
			expectedHour: 4,
			expectedLoc:  "Europe/Paris",
		},
		"invalid flag (1)": {
			flag:          "invalid",
			expectedError: "schedule-at=only contains two space-separated fields",
		},
		"invalid flag (2)": {
			flag:          "invalid invalid2",
			expectedError: "schedule-at=first field must be a digit representing the hour",
		},
		"invalid flag (3)": {
			flag:          "invalid invalid2:invalid3",
			expectedError: "schedule-at=first field must be a digit representing the hour",
		},
		"invalid hour": {
			flag:          "aa Europe/Paris",
			expectedError: "schedule-at=first field must be a digit representing the hour",
		},
		"unknown timezone": {
			flag:          "3 Europe/Unknown",
			expectedError: "schedule-at=unknown timezoneEurope/Unknown",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hour, loc, err := parseScheduleAtFlag(test.flag)
			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				assert.Equal(t, -1, hour)
				assert.Nil(t, loc)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, loc)
			assert.Equal(t, test.expectedHour, hour)
			assert.Equal(t, test.expectedLoc, loc.String())
		})
	}
}
