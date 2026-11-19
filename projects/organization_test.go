package projects

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/scalingomock"
)

type fakeProjectGetter struct {
	project          scalingo.Project
	err              error
	requestedProject []string
}

func (g *fakeProjectGetter) ProjectGet(_ context.Context, projectID string) (scalingo.Project, error) {
	g.requestedProject = append(g.requestedProject, projectID)
	return g.project, g.err
}

func TestOrganizationProjectSlug(t *testing.T) {
	acme := scalingo.Organization{ID: "org-acme", Slug: "acme", Name: "Acme Corp"}
	globex := scalingo.Organization{ID: "org-globex", Slug: "globex", Name: "Globex Inc"}

	tests := map[string]struct {
		projectID         string
		project           scalingo.Project
		projectErr        error
		organizations     []scalingo.Organization
		organizationsErr  error
		expectOrgsRead    bool
		expectedSlug      string
		expectedRequested []string
	}{
		"it should return the organization slug and the project name": {
			projectID:         "pr-1",
			project:           scalingo.Project{ID: "pr-1", Name: "test-project", OrganizationID: acme.ID},
			organizations:     []scalingo.Organization{acme},
			expectOrgsRead:    true,
			expectedSlug:      "acme/test-project",
			expectedRequested: []string{"pr-1"},
		},
		"it should match the organization on its exact ID": {
			projectID:         "pr-2",
			project:           scalingo.Project{ID: "pr-2", Name: "api", OrganizationID: globex.ID},
			organizations:     []scalingo.Organization{acme, globex},
			expectOrgsRead:    true,
			expectedSlug:      "globex/api",
			expectedRequested: []string{"pr-2"},
		},
		"it should not read the project nor the organizations without project ID": {
			projectID: "",
		},
		"it should not list the organizations for a personal project": {
			projectID:         "pr-3",
			project:           scalingo.Project{ID: "pr-3", Name: "default", Default: true},
			expectedRequested: []string{"pr-3"},
		},
		"it should return nothing when the organization is not found": {
			projectID:         "pr-4",
			project:           scalingo.Project{ID: "pr-4", Name: "app", OrganizationID: "org-unknown"},
			organizations:     []scalingo.Organization{acme, globex},
			expectOrgsRead:    true,
			expectedRequested: []string{"pr-4"},
		},
		"it should return nothing when the organization has no slug": {
			projectID:         "pr-5",
			project:           scalingo.Project{ID: "pr-5", Name: "app", OrganizationID: "org-noslug"},
			organizations:     []scalingo.Organization{{ID: "org-noslug"}},
			expectOrgsRead:    true,
			expectedRequested: []string{"pr-5"},
		},
		"it should return nothing when the project cannot be read": {
			projectID:         "pr-6",
			projectErr:        stderrors.New("forbidden"),
			expectedRequested: []string{"pr-6"},
		},
		"it should return nothing when the organizations cannot be listed": {
			projectID:         "pr-7",
			project:           scalingo.Project{ID: "pr-7", Name: "app", OrganizationID: acme.ID},
			organizationsErr:  stderrors.New("unavailable"),
			expectOrgsRead:    true,
			expectedRequested: []string{"pr-7"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			projectsService := &fakeProjectGetter{project: test.project, err: test.projectErr}
			organizationsService := scalingomock.NewMockOrganizationsService(gomock.NewController(t))
			if test.expectOrgsRead {
				organizationsService.EXPECT().OrganizationsList(ctx).Return(test.organizations, test.organizationsErr)
			}

			slug := OrganizationProjectSlug(ctx, projectsService, organizationsService, test.projectID)

			assert.Equal(t, test.expectedSlug, slug)
			assert.Equal(t, test.expectedRequested, projectsService.requestedProject)
		})
	}
}
