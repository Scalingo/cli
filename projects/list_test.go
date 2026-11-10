package projects

import (
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/scalingomock"
)

type fakeProjectsService struct {
	personal              []scalingo.Project
	personalErr           error
	byOrganization        map[string][]scalingo.Project
	byOrganizationErr     map[string]error
	requestedOrganization []string
	mutex                 *sync.Mutex
}

func (s *fakeProjectsService) ProjectsList(_ context.Context) ([]scalingo.Project, error) {
	return s.personal, s.personalErr
}

func (s *fakeProjectsService) ProjectsListByOrganization(_ context.Context, organizationID string) ([]scalingo.Project, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.requestedOrganization = append(s.requestedOrganization, organizationID)
	return s.byOrganization[organizationID], s.byOrganizationErr[organizationID]
}

func TestListWithOrganizations(t *testing.T) {
	acme := scalingo.Organization{ID: "org-acme", Slug: "acme", Name: "Acme Corp", Role: "admin"}
	globex := scalingo.Organization{ID: "org-globex", Slug: "globex", Name: "Globex Inc", Role: "member"}

	personalApp := scalingo.Project{ID: "pr-1", Name: "app", Default: true}
	acmeApp := scalingo.Project{ID: "pr-2", Name: "app", OrganizationID: acme.ID}
	acmeDB := scalingo.Project{ID: "pr-3", Name: "db", Default: true, OrganizationID: acme.ID}
	globexApp := scalingo.Project{ID: "pr-4", Name: "api", OrganizationID: globex.ID}

	// More organizations than concurrent requests: the order of the organizations must be kept.
	var manyOrganizations []scalingo.Organization
	manyProjects := map[string][]scalingo.Project{}
	var manyExpected []listedProject
	var manyRequested []string
	for i := range 3 * maxConcurrentRequests {
		organization := scalingo.Organization{ID: fmt.Sprintf("org-%d", i), Slug: fmt.Sprintf("slug-%d", i)}
		project := scalingo.Project{ID: fmt.Sprintf("pr-%d", i), Name: "app", OrganizationID: organization.ID}
		manyOrganizations = append(manyOrganizations, organization)
		manyProjects[organization.ID] = []scalingo.Project{project}
		manyExpected = append(manyExpected, listedProject{project: project, organizationSlug: organization.Slug})
		manyRequested = append(manyRequested, organization.ID)
	}

	tests := map[string]struct {
		projects          fakeProjectsService
		organizations     []scalingo.Organization
		organizationsErr  error
		expected          []listedProject
		expectedRequested []string
		expectedError     string
	}{
		"it should list the personal projects of a user without organization": {
			projects: fakeProjectsService{personal: []scalingo.Project{personalApp}},
			expected: []listedProject{{project: personalApp}},
		},
		"it should add the projects of an organization after the personal projects": {
			projects: fakeProjectsService{
				personal:       []scalingo.Project{personalApp},
				byOrganization: map[string][]scalingo.Project{acme.ID: {acmeApp, acmeDB}},
			},
			organizations: []scalingo.Organization{acme},
			expected: []listedProject{
				{project: personalApp},
				{project: acmeApp, organizationSlug: "acme"},
				{project: acmeDB, organizationSlug: "acme"},
			},
			expectedRequested: []string{acme.ID},
		},
		"it should associate the projects of several organizations with their own slug": {
			projects: fakeProjectsService{
				byOrganization: map[string][]scalingo.Project{acme.ID: {acmeApp}, globex.ID: {globexApp}},
			},
			organizations: []scalingo.Organization{acme, globex},
			expected: []listedProject{
				{project: acmeApp, organizationSlug: "acme"},
				{project: globexApp, organizationSlug: "globex"},
			},
			expectedRequested: []string{acme.ID, globex.ID},
		},
		"it should keep the order of the organizations when there are more than the concurrent requests": {
			projects:          fakeProjectsService{byOrganization: manyProjects},
			organizations:     manyOrganizations,
			expected:          manyExpected,
			expectedRequested: manyRequested,
		},
		"it should continue with the other organizations when one has no project": {
			projects: fakeProjectsService{
				byOrganization: map[string][]scalingo.Project{globex.ID: {globexApp}},
			},
			organizations:     []scalingo.Organization{acme, globex},
			expected:          []listedProject{{project: globexApp, organizationSlug: "globex"}},
			expectedRequested: []string{acme.ID, globex.ID},
		},
		"it should request the projects of an organization only once when it is listed twice": {
			projects: fakeProjectsService{
				byOrganization: map[string][]scalingo.Project{acme.ID: {acmeApp}},
			},
			organizations:     []scalingo.Organization{acme, acme},
			expected:          []listedProject{{project: acmeApp, organizationSlug: "acme"}},
			expectedRequested: []string{acme.ID},
		},
		"it should return no project when there is nothing to list": {
			organizations:     []scalingo.Organization{acme},
			expected:          []listedProject{},
			expectedRequested: []string{acme.ID},
		},
		"it should return an error when the personal projects cannot be listed": {
			projects:      fakeProjectsService{personalErr: stderrors.New("boom")},
			expectedError: "list projects",
		},
		"it should return an error when the organizations cannot be listed": {
			projects:         fakeProjectsService{personal: []scalingo.Project{personalApp}},
			organizationsErr: stderrors.New("boom"),
			expectedError:    "list organizations",
		},
		"it should return an error naming the organization when its projects cannot be listed": {
			projects: fakeProjectsService{
				personal:          []scalingo.Project{personalApp},
				byOrganization:    map[string][]scalingo.Project{acme.ID: {acmeApp}},
				byOrganizationErr: map[string]error{globex.ID: stderrors.New("forbidden")},
			},
			organizations:     []scalingo.Organization{acme, globex},
			expectedError:     "list projects of organization globex",
			expectedRequested: []string{acme.ID, globex.ID},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			test.projects.mutex = &sync.Mutex{}
			organizationsService := scalingomock.NewMockOrganizationsService(gomock.NewController(t))
			if test.projects.personalErr == nil {
				organizationsService.EXPECT().OrganizationsList(ctx).Return(test.organizations, test.organizationsErr)
			}

			projects, err := listWithOrganizations(ctx, &test.projects, organizationsService)

			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				assert.Nil(t, projects)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, projects)
			}
			assert.ElementsMatch(t, test.expectedRequested, test.projects.requestedOrganization)
		})
	}
}

func TestRender(t *testing.T) {
	t.Run("it should render the name, organization slug, default and ID of the projects", func(t *testing.T) {
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
