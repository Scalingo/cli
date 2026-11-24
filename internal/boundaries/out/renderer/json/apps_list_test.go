package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-scalingo/v11"
)

func TestAppsListRenderer_Render(t *testing.T) {
	personalApp := &scalingo.App{ID: "app-1", Name: "personal-web", Status: scalingo.AppStatusRunning}
	organizationApp := &scalingo.App{ID: "app-2", Name: "acme-web", Status: scalingo.AppStatusStopped}

	tests := map[string]struct {
		listedApps []apps.ListedApp
		expected   string
	}{
		"it renders the apps of all the scopes in the apps envelope": {
			listedApps: []apps.ListedApp{
				{App: personalApp},
				{App: organizationApp, OrganizationID: "org-acme", OrganizationSlug: "acme"},
			},
			expected: `{"apps":[` + mustMarshal(t, personalApp) + `,` + mustMarshal(t, organizationApp) + `]}` + "\n",
		},
		"it renders an empty list for an empty list of apps": {
			listedApps: []apps.ListedApp{},
			expected:   `{"apps":[]}` + "\n",
		},
		"it renders a null list for a nil list of apps": {
			listedApps: nil,
			expected:   `{"apps":null}` + "\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewAppsList()
			r.SetData(t.Context(), test.listedApps)

			output := captureStdout(t, func() {
				require.NoError(t, r.Render(t.Context()))
			})

			assert.Equal(t, test.expected, output)
			assert.NotContains(t, output, "org-acme")
		})
	}
}

func mustMarshal(t *testing.T, app *scalingo.App) string {
	t.Helper()

	data, err := json.Marshal(app)
	require.NoError(t, err)
	return string(data)
}
