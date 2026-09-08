package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHoursBeforeDeleteValidator(t *testing.T) {
	var hoursBeforeDelete uint
	validate := hoursBeforeDeleteValidator(t.Context(), &hoursBeforeDelete)

	for _, value := range []string{"", "0", "1", "2147483647"} {
		t.Run("accepts "+value, func(t *testing.T) {
			require.NoError(t, validate(value))
		})
	}

	tests := map[string]struct {
		value string
		error string
	}{
		"negative":     {value: "-1", error: "must be positive"},
		"not a number": {value: "later", error: "error parsing hours"},
		"out of range": {value: "2147483648", error: "error parsing hours"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.ErrorContains(t, validate(test.value), test.error)
		})
	}
}
