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

func TestOrganizationProjectSlug(t *testing.T) {
	acme := scalingo.Organization{ID: "org-acme", Slug: "acme", Name: "Acme Corp"}
	globex := scalingo.Organization{ID: "org-globex", Slug: "globex", Name: "Globex Inc"}

	tests := map[string]struct {
		projectID    string
		expectations func(ctx context.Context, client *scalingomock.MockAPI)
		expectedSlug string
	}{
		"it should return the organization slug and the project name": {
			projectID: "pr-1",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectGet(ctx, "pr-1").Return(scalingo.Project{ID: "pr-1", Name: "test-project", OrganizationID: acme.ID}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme}, nil)
			},
			expectedSlug: "acme/test-project",
		},
		"it should match the organization on its exact ID": {
			projectID: "pr-2",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectGet(ctx, "pr-2").Return(scalingo.Project{ID: "pr-2", Name: "api", OrganizationID: globex.ID}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
			},
			expectedSlug: "globex/api",
		},
		"it should not read the project nor the organizations without project ID": {
			projectID:    "",
			expectations: func(_ context.Context, _ *scalingomock.MockAPI) {},
		},
		"it should not list the organizations for a personal project": {
			projectID: "pr-3",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectGet(ctx, "pr-3").Return(scalingo.Project{ID: "pr-3", Name: "default", Default: true}, nil)
			},
		},
		"it should return nothing when the organization is not found": {
			projectID: "pr-4",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectGet(ctx, "pr-4").Return(scalingo.Project{ID: "pr-4", Name: "app", OrganizationID: "org-unknown"}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
			},
		},
		"it should return nothing when the organization has no slug": {
			projectID: "pr-5",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectGet(ctx, "pr-5").Return(scalingo.Project{ID: "pr-5", Name: "app", OrganizationID: "org-noslug"}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{{ID: "org-noslug"}}, nil)
			},
		},
		"it should return nothing when the project cannot be read": {
			projectID: "pr-6",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectGet(ctx, "pr-6").Return(scalingo.Project{}, stderrors.New("forbidden"))
			},
		},
		"it should return nothing when the organizations cannot be listed": {
			projectID: "pr-7",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectGet(ctx, "pr-7").Return(scalingo.Project{ID: "pr-7", Name: "app", OrganizationID: acme.ID}, nil)
				client.EXPECT().OrganizationsList(ctx).Return(nil, stderrors.New("unavailable"))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			client := scalingomock.NewMockAPI(gomock.NewController(t))
			test.expectations(ctx, client)

			slug := OrganizationProjectSlug(ctx, client, test.projectID)

			assert.Equal(t, test.expectedSlug, slug)
		})
	}
}
