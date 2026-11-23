package dbdr

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-scalingo/v11"
)

func TestRender(t *testing.T) {
	const organizationSlug = "acme"
	currentUser := &scalingo.User{Email: "john@example.com"}
	newDatabase := func(id, projectName, ownerEmail string) scalingo.DatabaseNG {
		database := scalingo.DatabaseNG{
			ID:         id,
			Name:       id + "-name",
			Technology: "postgresql",
			Plan:       "postgresql-dr-starter-4096",
			Database:   scalingo.Database{Status: scalingo.DatabaseStatusRunning},
			App:        scalingo.App{Owner: scalingo.Owner{Email: ownerEmail}},
		}
		database.App.Project.Name = projectName
		return database
	}

	t.Run("it should distinguish projects with the same name in the personal scope and in an organization", func(t *testing.T) {
		output := &bytes.Buffer{}

		err := render(output, currentUser, []listedDatabase{
			{database: newDatabase("db-1", "default", currentUser.Email)},
			{database: newDatabase("db-2", "default", "jane@example.com"), organizationSlug: organizationSlug},
		})

		require.NoError(t, err)
		assert.Regexp(t, `│\s+ID\s+│\s+NAME\s+│\s+TYPE\s+│\s+PLAN\s+│\s+ROLE\s+│\s+STATUS\s+│\s+PROJECT\s+│`, output.String())
		assert.Regexp(t, `│\s+db-1\s+│\s+db-1-name\s+│\s+postgresql\s+│\s+postgresql-dr-starter-4096\s+│\s+owner\s+│\s+running\s+│\s+default\s+│`, output.String())
		assert.Regexp(t, `│\s+db-2\s+│\s+db-2-name\s+│\s+postgresql\s+│\s+postgresql-dr-starter-4096\s+│\s+collaborator\s+│\s+running\s+│\s+acme/default\s+│`, output.String())
	})

	t.Run("it should render an empty project cell when the project name is unknown", func(t *testing.T) {
		output := &bytes.Buffer{}

		err := render(output, currentUser, []listedDatabase{
			{database: newDatabase("db-1", "", currentUser.Email), organizationSlug: organizationSlug},
		})

		require.NoError(t, err)
		assert.Regexp(t, `│\s+running\s+│\s+│`, output.String())
		assert.NotContains(t, output.String(), organizationSlug)
	})

	t.Run("it should only render the headers when there is no database", func(t *testing.T) {
		output := &bytes.Buffer{}

		err := render(output, currentUser, nil)

		require.NoError(t, err)
		assert.Regexp(t, `│\s+ID\s+│\s+NAME\s+│\s+TYPE\s+│\s+PLAN\s+│\s+ROLE\s+│\s+STATUS\s+│\s+PROJECT\s+│`, output.String())
	})
}
