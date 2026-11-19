package projects

import (
	"context"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// The SDK builds its requests without the context. The context is therefore only a budget checked between the two reads.
	t.Run("it should not list the organizations when the budget is exhausted by the project read", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		projectsService := &cancelingProjectGetter{cancel: cancel, project: scalingo.Project{ID: "pr-1", Name: "app", OrganizationID: acme.ID}}
		organizationsService := scalingomock.NewMockOrganizationsService(gomock.NewController(t))

		slug := OrganizationProjectSlug(ctx, projectsService, organizationsService, "pr-1")

		assert.Empty(t, slug)
	})
}

type cancelingProjectGetter struct {
	cancel  context.CancelFunc
	project scalingo.Project
}

func (g *cancelingProjectGetter) ProjectGet(_ context.Context, _ string) (scalingo.Project, error) {
	g.cancel()
	return g.project, nil
}

// fakeAPI serves the project and organizations endpoints used by the enrichment. A handler named in `slow` blocks
// until the client gives up, which is what a stalled API looks like to the CLI.
type fakeAPI struct {
	projectRequests      atomic.Int32
	organizationRequests atomic.Int32
}

func (a *fakeAPI) start(t *testing.T, project string, slow ...string) *httptest.Server {
	t.Helper()

	isSlow := func(name string) bool {
		for _, s := range slow {
			if s == name {
				return true
			}
		}
		return false
	}
	block := func(r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(10 * time.Second):
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/projects/pr-1", func(w http.ResponseWriter, r *http.Request) {
		a.projectRequests.Add(1)
		if isSlow("project") {
			block(r)
			return
		}
		_, _ = w.Write([]byte(project))
	})
	mux.HandleFunc("/v1/organizations", func(w http.ResponseWriter, r *http.Request) {
		a.organizationRequests.Add(1)
		if isSlow("organizations") {
			block(r)
			return
		}
		_, _ = w.Write([]byte(`{"organizations":[{"id":"org-acme","slug":"acme"}]}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

// newTestClient returns a client wired like the CLI one, except for the token which is static.
func newTestClient(server *httptest.Server) clientFactory {
	return func(ctx context.Context, timeout time.Duration) (*scalingo.Client, error) {
		return scalingo.New(ctx, scalingo.ClientConfig{
			APIEndpoint:          server.URL,
			AuthEndpoint:         server.URL,
			Timeout:              timeout,
			StaticTokenGenerator: scalingo.NewStaticTokenGenerator("token"),
		})
	}
}

func TestCreatedProjectSlug(t *testing.T) {
	const orgProject = `{"project":{"id":"pr-1","name":"test-project","organization_id":"org-acme"}}`
	const timeout = 200 * time.Millisecond
	// A stalled request ends with the timeout. The handlers would otherwise block for 10 seconds.
	const maxDuration = 3 * time.Second

	t.Run("it should resolve the slug through the project and organizations endpoints", func(t *testing.T) {
		api := &fakeAPI{}
		server := api.start(t, orgProject)

		slug := createdProjectSlug(t.Context(), newTestClient(server), timeout, "pr-1")

		assert.Equal(t, "acme/test-project", slug)
		assert.EqualValues(t, 1, api.projectRequests.Load())
		assert.EqualValues(t, 1, api.organizationRequests.Load())
	})

	t.Run("it should give up on a stalled project read and not list the organizations", func(t *testing.T) {
		api := &fakeAPI{}
		server := api.start(t, orgProject, "project")

		start := time.Now()
		slug := createdProjectSlug(t.Context(), newTestClient(server), timeout, "pr-1")

		assert.Empty(t, slug)
		assert.Less(t, time.Since(start), maxDuration)
		assert.EqualValues(t, 1, api.projectRequests.Load())
		assert.EqualValues(t, 0, api.organizationRequests.Load())
	})

	t.Run("it should give up on a stalled organizations read", func(t *testing.T) {
		api := &fakeAPI{}
		server := api.start(t, orgProject, "organizations")

		start := time.Now()
		slug := createdProjectSlug(t.Context(), newTestClient(server), timeout, "pr-1")

		assert.Empty(t, slug)
		assert.Less(t, time.Since(start), maxDuration)
		assert.EqualValues(t, 1, api.projectRequests.Load())
		assert.EqualValues(t, 1, api.organizationRequests.Load())
	})

	t.Run("it should not create a client without project ID", func(t *testing.T) {
		newClient := func(_ context.Context, _ time.Duration) (*scalingo.Client, error) {
			t.Error("the client must not be created without project ID")
			return nil, nil
		}

		assert.Empty(t, createdProjectSlug(t.Context(), newClient, timeout, ""))
	})

	t.Run("it should return nothing when the client cannot be created", func(t *testing.T) {
		newClient := func(_ context.Context, _ time.Duration) (*scalingo.Client, error) {
			return nil, stderrors.New("unauthenticated")
		}

		assert.Empty(t, createdProjectSlug(t.Context(), newClient, timeout, "pr-1"))
	})

	t.Run("it should request the client with the given timeout", func(t *testing.T) {
		var requested time.Duration
		newClient := func(ctx context.Context, timeout time.Duration) (*scalingo.Client, error) {
			requested = timeout
			return nil, stderrors.New("stop here")
		}

		createdProjectSlug(t.Context(), newClient, 1234*time.Millisecond, "pr-1")

		assert.Equal(t, 1234*time.Millisecond, requested)
	})

	t.Run("it should use a short timeout for the creation flows", func(t *testing.T) {
		require.LessOrEqual(t, enrichmentTimeout, 3*time.Second)
	})
}
