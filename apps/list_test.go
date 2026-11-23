package apps

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/go-scalingo/v11"
)

func TestFilterAppsByProject(t *testing.T) {
	const (
		organizationID   = "org-acme"
		organizationSlug = "acme"
		projectName      = "app"
	)
	owner := scalingo.Owner{Username: "john"}
	personalApp := &scalingo.App{Name: "personal-web", Owner: owner}
	personalApp.Project.Name = projectName
	organizationApp := &scalingo.App{Name: "acme-web", Owner: owner}
	organizationApp.Project.Name = projectName
	otherProjectApp := &scalingo.App{Name: "acme-worker", Owner: owner}
	otherProjectApp.Project.Name = "back"

	apps := []ListedApp{
		{App: personalApp},
		{App: organizationApp, OrganizationID: organizationID, OrganizationSlug: organizationSlug},
		{App: otherProjectApp, OrganizationID: organizationID, OrganizationSlug: organizationSlug},
	}

	tests := map[string]struct {
		projectSlug string
		expected    []ListedApp
	}{
		"it should keep all the apps without filter": {
			expected: apps,
		},
		"it should filter the apps of all the scopes with the <ownerUsername>/<projectName> format": {
			projectSlug: owner.Username + "/" + projectName,
			expected:    []ListedApp{apps[0], apps[1]},
		},
		"it should not filter with the organization slug displayed in the table": {
			projectSlug: organizationSlug + "/" + projectName,
			expected:    []ListedApp{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, filterAppsByProject(apps, test.projectSlug))
		})
	}
}
