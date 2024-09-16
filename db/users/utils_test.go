package users_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/cli/db/users"
)

func Test_IsPasswordValid(t *testing.T) {
	testPasswords := map[string]struct {
		password         string
		confirmation     string
		expectedValidity bool
		expectedMessage  string
	}{
		"empty": {
			password:         "",
			confirmation:     "",
			expectedValidity: true,
			expectedMessage:  "",
		},
		"confirmation doesn't match": {
			password:         "abc",
			confirmation:     "aBc",
			expectedValidity: false,
			expectedMessage:  "Password confirmation doesn't match",
		},
		"too short": {
			password:         "123456789a123456789b123",
			confirmation:     "123456789a123456789b123",
			expectedValidity: false,
			expectedMessage:  "Password must contain between 24 and 64 characters",
		},
		"too long": {
			password:         "123456789a123456789b123456789c123456789d123456789e123456789f12345",
			confirmation:     "123456789a123456789b123456789c123456789d123456789e123456789f12345",
			expectedValidity: false,
			expectedMessage:  "Password must contain between 24 and 64 characters",
		},
		"valid, short password": {
			password:         "123456789a123456789b1234",
			confirmation:     "123456789a123456789b1234",
			expectedValidity: true,
			expectedMessage:  "",
		},
		"valid, log password ": {
			password:         "123456789a123456789b123456789c123456789d123456789e123456789f1234",
			confirmation:     "123456789a123456789b123456789c123456789d123456789e123456789f1234",
			expectedValidity: true,
			expectedMessage:  "",
		},
	}

	for name, testCase := range testPasswords {
		t.Run(name, func(t *testing.T) {
			message, isValid := users.IsPasswordValid(testCase.password, testCase.confirmation)

			assert.Equal(t, testCase.expectedValidity, isValid)
			assert.Equal(t, testCase.expectedMessage, message)
		})
	}
}

func Test_IsUsernameValid(t *testing.T) {
	testPasswords := map[string]struct {
		username         string
		expectedValidity bool
		expectedMessage  string
	}{
		"empty": {
			username:         "",
			expectedValidity: false,
			expectedMessage:  "Name must contain between 6 and 32 characters",
		},
		"too short": {
			username:         "12345",
			expectedValidity: false,
			expectedMessage:  "Name must contain between 6 and 32 characters",
		},
		"too long": {
			username:         "123456789a123456789b123456789c123",
			expectedValidity: false,
			expectedMessage:  "Name must contain between 6 and 32 characters",
		},
		"valid, short username": {
			username:         "123456",
			expectedValidity: true,
			expectedMessage:  "",
		},
		"valid, long username": {
			username:         "123456789a123456789b123456789c12",
			expectedValidity: true,
			expectedMessage:  "",
		},
	}

	for name, testCase := range testPasswords {
		t.Run(name, func(t *testing.T) {
			message, isValid := users.IsUsernameValid(testCase.username)

			assert.Equal(t, testCase.expectedValidity, isValid)
			assert.Equal(t, testCase.expectedMessage, message)
		})
	}
}
