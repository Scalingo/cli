package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationProjectLabel(t *testing.T) {
	const (
		organizationSlug = "acme"
		projectName      = "app"
	)

	tests := map[string]struct {
		organizationSlug string
		projectName      string
		expected         string
	}{
		"it prefixes the project name with the organization slug": {
			organizationSlug: organizationSlug,
			projectName:      projectName,
			expected:         "acme/app",
		},
		"it returns only the project name when the organization slug is unknown": {
			projectName: projectName,
			expected:    projectName,
		},
		"it returns an empty label when the project name is unknown": {
			organizationSlug: organizationSlug,
			expected:         "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, OrganizationProjectLabel(test.organizationSlug, test.projectName))
		})
	}
}
