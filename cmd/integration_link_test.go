package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateHoursBeforeDelete(t *testing.T) {
	validate := validateHoursBeforeDelete(t.Context())

	for _, value := range []string{"", "0", "1", "2147483647"} {
		t.Run("accepts "+value, func(t *testing.T) {
			require.NoError(t, validate(value))
		})
	}

	for _, test := range []struct {
		name  string
		value string
		error string
	}{
		{name: "negative", value: "-1", error: "must be positive"},
		{name: "not a number", value: "later", error: "error parsing hours"},
		{name: "out of range", value: "2147483648", error: "error parsing hours"},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.ErrorContains(t, validate(test.value), test.error)
		})
	}
}

func TestParseHoursBeforeDelete(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint
	}{
		{name: "empty uses default", input: "", expected: 0},
		{name: "zero", input: "0", expected: 0},
		{name: "positive", input: "12", expected: 12},
		{name: "maximum", input: "2147483647", expected: 2147483647},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseHoursBeforeDelete(t.Context(), test.input)
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}
