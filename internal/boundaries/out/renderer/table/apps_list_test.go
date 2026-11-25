package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-scalingo/v11"
)

func TestAppsListRenderer_Render(t *testing.T) {
	const (
		organizationID = "org-acme"
		projectName    = "app"
	)
	acme := &scalingo.Organization{ID: organizationID, Slug: "acme"}
	currentUser := &scalingo.User{Email: "john@example.com"}
	owner := scalingo.Owner{Username: "john", Email: currentUser.Email}
	newApp := func(name, projectName string, status scalingo.AppStatus, owner scalingo.Owner) *scalingo.App {
		app := &scalingo.App{Name: name, Status: status, Owner: owner}
		app.Project.Name = projectName
		return app
	}

	t.Run("it distinguishes projects with the same name in the personal scope and in an organization", func(t *testing.T) {
		r := NewAppsList(currentUser)
		r.SetData(t.Context(), []apps.ListedApp{
			{App: newApp("personal-web", projectName, scalingo.AppStatusRunning, owner)},
			{
				App:          newApp("acme-web", projectName, scalingo.AppStatusStopped, scalingo.Owner{Username: "jane", Email: "jane@example.com"}),
				Organization: acme,
			},
		})

		output := captureStdout(t, func() {
			require.NoError(t, r.Render(t.Context()))
		})

		assert.Regexp(t, `│\s+NAME\s+│\s+ROLE\s+│\s+STATUS\s+│\s+PROJECT\s+│`, output)
		assert.Regexp(t, `│\s+personal-web\s+│\s+owner\s+│\s+running\s+│\s+john/app\s+│`, output)
		assert.Regexp(t, `│\s+acme-web\s+│\s+collaborator\s+│\s+stopped\s+│\s+acme/app\s+│`, output)
		assert.NotContains(t, output, organizationID)
	})

	t.Run("it displays only the project name when the organization slug is unknown", func(t *testing.T) {
		r := NewAppsList(currentUser)
		r.SetData(t.Context(), []apps.ListedApp{
			{App: newApp("no-slug", projectName, scalingo.AppStatusRunning, owner), Organization: &scalingo.Organization{ID: organizationID}},
			{App: newApp("no-project", "", scalingo.AppStatusRunning, owner), Organization: acme},
		})

		output := captureStdout(t, func() {
			require.NoError(t, r.Render(t.Context()))
		})

		assert.Regexp(t, `│\s+no-slug\s+│\s+owner\s+│\s+running\s+│\s+app\s+│`, output)
		assert.Regexp(t, `│\s+no-project\s+│\s+owner\s+│\s+running\s+│\s+│`, output)
		assert.NotContains(t, output, "john/app")
		assert.NotContains(t, output, "acme/")
	})

	t.Run("it displays a message when there is no app", func(t *testing.T) {
		r := NewAppsList(currentUser)
		r.SetData(t.Context(), []apps.ListedApp{})

		output := captureStdout(t, func() {
			require.NoError(t, r.Render(t.Context()))
		})

		assert.Contains(t, output, "You haven't created any app yet")
	})
}
