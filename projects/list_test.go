package projects

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
	acme := scalingo.Organization{ID: "org-acme", Slug: "acme", Name: "Acme Corp", Role: "admin"}
	globex := scalingo.Organization{ID: "org-globex", Slug: "globex", Name: "Globex Inc", Role: "member"}

	personalApp := scalingo.Project{ID: "pr-1", Name: "app", Default: true}
	acmeApp := scalingo.Project{ID: "pr-2", Name: "app", OrganizationID: acme.ID}
	acmeDB := scalingo.Project{ID: "pr-3", Name: "db", Default: true, OrganizationID: acme.ID}
	globexApp := scalingo.Project{ID: "pr-4", Name: "api", OrganizationID: globex.ID}

	tests := map[string]struct {
		expectations  func(ctx context.Context, client *scalingomock.MockAPI)
		expected      []listedProject
		expectedError string
	}{
		"it should list the personal projects of a user without organization": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return([]scalingo.Project{personalApp}, nil)
				client.EXPECT().OrganizationsList(ctx).Return(nil, nil)
			},
			expected: []listedProject{{project: personalApp}},
		},
		"it should add the projects of an organization after the personal projects": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return([]scalingo.Project{personalApp}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme}, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, acme.ID).Return([]scalingo.Project{acmeApp, acmeDB}, nil)
			},
			expected: []listedProject{
				{project: personalApp},
				{project: acmeApp, organizationSlug: "acme"},
				{project: acmeDB, organizationSlug: "acme"},
			},
		},
		"it should associate the projects of several organizations with their own slug": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return(nil, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, acme.ID).Return([]scalingo.Project{acmeApp}, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, globex.ID).Return([]scalingo.Project{globexApp}, nil)
			},
			expected: []listedProject{
				{project: acmeApp, organizationSlug: "acme"},
				{project: globexApp, organizationSlug: "globex"},
			},
		},
		"it should continue with the other organizations when one has no project": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return(nil, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, acme.ID).Return(nil, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, globex.ID).Return([]scalingo.Project{globexApp}, nil)
			},
			expected: []listedProject{{project: globexApp, organizationSlug: "globex"}},
		},
		"it should return no project when there is nothing to list": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return(nil, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme}, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, acme.ID).Return(nil, nil)
			},
			expected: []listedProject{},
		},
		"it should return an error when the personal projects cannot be listed": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return(nil, stderrors.New("boom"))
			},
			expectedError: "list projects",
		},
		"it should return an error when the organizations cannot be listed": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return([]scalingo.Project{personalApp}, nil)
				client.EXPECT().OrganizationsList(ctx).Return(nil, stderrors.New("boom"))
			},
			expectedError: "list organizations",
		},
		"it should return an error when the projects of an organization cannot be listed": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectsList(ctx).Return([]scalingo.Project{personalApp}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, acme.ID).Return([]scalingo.Project{acmeApp}, nil)
				client.EXPECT().ProjectsListByOrganization(ctx, globex.ID).Return(nil, stderrors.New("forbidden"))
			},
			expectedError: "list organization projects",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			client := scalingomock.NewMockAPI(gomock.NewController(t))
			test.expectations(ctx, client)

			projects, err := listWithOrganizations(ctx, client)

			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				assert.Nil(t, projects)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, projects)
			}
		})
	}
}

func TestRender(t *testing.T) {
	t.Run("it should render the name, organization, default and ID of the projects", func(t *testing.T) {
		output := &bytes.Buffer{}

		err := render(output, []listedProject{
			{project: scalingo.Project{ID: "pr-1", Name: "app", Default: true}},
			{project: scalingo.Project{ID: "pr-2", Name: "app", OrganizationID: "org-uuid-1"}, organizationSlug: "acme"},
		})

		require.NoError(t, err)
		assert.Regexp(t, `│\s+NAME\s+│\s+ORGANIZATION\s+│\s+DEFAULT\s+│\s+ID\s+│`, output.String())
		assert.Regexp(t, `│\s+app\s+│\s+│\s+true\s+│\s+pr-1\s+│`, output.String())
		assert.Regexp(t, `│\s+app\s+│\s+acme\s+│\s+false\s+│\s+pr-2\s+│`, output.String())
		assert.NotContains(t, output.String(), "org-uuid-1")
	})

	t.Run("it should only render the headers when there is no project", func(t *testing.T) {
		output := &bytes.Buffer{}

		err := render(output, nil)

		require.NoError(t, err)
		assert.Regexp(t, `│\s+NAME\s+│\s+ORGANIZATION\s+│\s+DEFAULT\s+│\s+ID\s+│`, output.String())
	})
}
