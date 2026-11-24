package dbdr

import (
	"bytes"
	"context"
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/scalingomock"
)

func TestListWithOrganizations(t *testing.T) {
	errForbidden := stderrors.New("forbidden")
	acme := scalingo.Organization{ID: "org-acme", Slug: "acme", Name: "Acme Corp", Role: "admin"}
	globex := scalingo.Organization{ID: "org-globex", Slug: "globex", Name: "Globex Inc", Role: "member"}

	personalDB := scalingo.DatabaseNG{ID: "db-1", Name: "personal-db"}
	acmeDB := scalingo.DatabaseNG{ID: "db-2", Name: "acme-db"}
	globexDB := scalingo.DatabaseNG{ID: "db-3", Name: "globex-db"}

	tests := map[string]struct {
		expectations  func(ctx context.Context, organizations *scalingomock.MockOrganizationsService, databases *scalingomock.MockDatabasesPreviewService)
		expected      []listedDatabase
		expectedError string
		expectedCause error
	}{
		"it lists only the personal databases of a user without organization": {
			expectations: func(ctx context.Context, organizations *scalingomock.MockOrganizationsService, databases *scalingomock.MockDatabasesPreviewService) {
				databases.EXPECT().DatabasesList(ctx).Return([]scalingo.DatabaseNG{personalDB}, nil)
				organizations.EXPECT().OrganizationsList(ctx).Return(nil, nil)
			},
			expected: []listedDatabase{{database: personalDB}},
		},
		"it adds the databases of each organization after the personal databases": {
			expectations: func(ctx context.Context, organizations *scalingomock.MockOrganizationsService, databases *scalingomock.MockDatabasesPreviewService) {
				databases.EXPECT().DatabasesList(ctx).Return([]scalingo.DatabaseNG{personalDB}, nil)
				organizations.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				databases.EXPECT().DatabasesListByOrganization(ctx, acme.ID).Return([]scalingo.DatabaseNG{acmeDB}, nil)
				databases.EXPECT().DatabasesListByOrganization(ctx, globex.ID).Return([]scalingo.DatabaseNG{globexDB}, nil)
			},
			expected: []listedDatabase{
				{database: personalDB},
				{database: acmeDB, organizationSlug: acme.Slug},
				{database: globexDB, organizationSlug: globex.Slug},
			},
		},
		"it still queries an organization without database": {
			expectations: func(ctx context.Context, organizations *scalingomock.MockOrganizationsService, databases *scalingomock.MockDatabasesPreviewService) {
				databases.EXPECT().DatabasesList(ctx).Return(nil, nil)
				organizations.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				databases.EXPECT().DatabasesListByOrganization(ctx, acme.ID).Return([]scalingo.DatabaseNG{}, nil)
				databases.EXPECT().DatabasesListByOrganization(ctx, globex.ID).Return([]scalingo.DatabaseNG{globexDB}, nil)
			},
			expected: []listedDatabase{{database: globexDB, organizationSlug: globex.Slug}},
		},
		"it returns an error when the personal databases cannot be listed": {
			expectations: func(ctx context.Context, _ *scalingomock.MockOrganizationsService, databases *scalingomock.MockDatabasesPreviewService) {
				databases.EXPECT().DatabasesList(ctx).Return(nil, stderrors.New("boom"))
			},
			expectedError: "list databases",
		},
		"it returns an error when the organizations cannot be listed": {
			expectations: func(ctx context.Context, organizations *scalingomock.MockOrganizationsService, databases *scalingomock.MockDatabasesPreviewService) {
				databases.EXPECT().DatabasesList(ctx).Return([]scalingo.DatabaseNG{personalDB}, nil)
				organizations.EXPECT().OrganizationsList(ctx).Return(nil, stderrors.New("boom"))
			},
			expectedError: "list organizations",
		},
		"it returns an error when organization databases cannot be listed": {
			expectations: func(ctx context.Context, organizations *scalingomock.MockOrganizationsService, databases *scalingomock.MockDatabasesPreviewService) {
				databases.EXPECT().DatabasesList(ctx).Return([]scalingo.DatabaseNG{personalDB}, nil)
				organizations.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				databases.EXPECT().DatabasesListByOrganization(ctx, acme.ID).Return([]scalingo.DatabaseNG{acmeDB}, nil)
				databases.EXPECT().DatabasesListByOrganization(ctx, globex.ID).Return(nil, errForbidden)
			},
			expectedError: "list organization databases",
			expectedCause: errForbidden,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			ctrl := gomock.NewController(t)
			organizations := scalingomock.NewMockOrganizationsService(ctrl)
			databases := scalingomock.NewMockDatabasesPreviewService(ctrl)
			test.expectations(ctx, organizations, databases)

			listed, err := listWithOrganizations(ctx, organizations, databases)

			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				if test.expectedCause != nil {
					require.ErrorIs(t, err, test.expectedCause)
				}
				assert.Nil(t, listed)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, listed)
			}
		})
	}
}

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

	t.Run("it distinguishes projects with the same name in the personal scope and in an organization", func(t *testing.T) {
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

	t.Run("it renders an empty project cell when the project name is unknown", func(t *testing.T) {
		output := &bytes.Buffer{}

		err := render(output, currentUser, []listedDatabase{
			{database: newDatabase("db-1", "", currentUser.Email), organizationSlug: organizationSlug},
		})

		require.NoError(t, err)
		assert.Regexp(t, `│\s+running\s+│\s+│`, output.String())
		assert.NotContains(t, output.String(), organizationSlug)
	})

	t.Run("it renders only the headers when there is no database", func(t *testing.T) {
		output := &bytes.Buffer{}

		err := render(output, currentUser, nil)

		require.NoError(t, err)
		assert.Regexp(t, `│\s+ID\s+│\s+NAME\s+│\s+TYPE\s+│\s+PLAN\s+│\s+ROLE\s+│\s+STATUS\s+│\s+PROJECT\s+│`, output.String())
	})
}
