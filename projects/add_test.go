package projects

import (
	"context"
	stderrors "errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-scalingo/v11"
	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-scalingo/v11/scalingomock"
)

func TestCreateProject(t *testing.T) {
	other := scalingo.Organization{ID: "org-other", Slug: "other-org", Name: "Other", Role: "owner"}
	myOrg := scalingo.Organization{ID: "org-11111111-2222-4333-8444-555555555555", Slug: "my-org-1", Name: "My Org", Role: "admin"}

	authErr := stderrors.New("auth unavailable")
	forbiddenErr := &httpclient.RequestFailedError{
		Code:     http.StatusForbidden,
		APIError: httpclient.ForbiddenError{Err: "you are not allowed to create a project in this organization"},
	}

	tests := map[string]struct {
		params             scalingo.ProjectAddParams
		organizationSlug   string
		expectations       func(ctx context.Context, client *scalingomock.MockAPI)
		expected           scalingo.Project
		expectedError      string
		expectedErrorCause error
	}{
		"it creates the project in the organization matching the slug": {
			params:           scalingo.ProjectAddParams{Name: "my-project"},
			organizationSlug: myOrg.Slug,
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				gomock.InOrder(
					client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{other, myOrg}, nil),
					client.EXPECT().ProjectAdd(ctx, scalingo.ProjectAddParams{Name: "my-project", OrganizationID: myOrg.ID}).
						Return(scalingo.Project{ID: "prj-1", Name: "my-project-from-api", OrganizationID: myOrg.ID}, nil),
				)
			},
			expected: scalingo.Project{ID: "prj-1", Name: "my-project-from-api", OrganizationID: myOrg.ID},
		},
		"it creates a personal project without listing the organizations": {
			params: scalingo.ProjectAddParams{Name: "my-project", Default: true},
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().ProjectAdd(ctx, scalingo.ProjectAddParams{Name: "my-project", Default: true}).
					Return(scalingo.Project{ID: "prj-2", Name: "my-project", Default: true}, nil)
			},
			expected: scalingo.Project{ID: "prj-2", Name: "my-project", Default: true},
		},
		"it does not create the project when no organization matches the slug": {
			params:           scalingo.ProjectAddParams{Name: "my-project"},
			organizationSlug: "unknown-org",
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{other, myOrg}, nil)
			},
			expectedError: `organization "unknown-org" not found`,
		},
		"it does not match an organization by its ID": {
			params:           scalingo.ProjectAddParams{Name: "my-project"},
			organizationSlug: myOrg.ID,
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{other, myOrg}, nil)
			},
			expectedError: `organization "` + myOrg.ID + `" not found`,
		},
		"it does not create the project when the matching organization has no ID": {
			params:           scalingo.ProjectAddParams{Name: "my-project"},
			organizationSlug: myOrg.Slug,
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{{Slug: myOrg.Slug}}, nil)
			},
			expectedError: `organization "my-org-1" not found`,
		},
		"it does not create the project when the organizations cannot be listed": {
			params:           scalingo.ProjectAddParams{Name: "my-project"},
			organizationSlug: myOrg.Slug,
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().OrganizationsList(ctx).Return(nil, authErr)
			},
			expectedError:      "list organizations",
			expectedErrorCause: authErr,
		},
		"it returns the API error without retrying": {
			params:           scalingo.ProjectAddParams{Name: "my-project"},
			organizationSlug: myOrg.Slug,
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{other, myOrg}, nil)
				client.EXPECT().ProjectAdd(ctx, scalingo.ProjectAddParams{Name: "my-project", OrganizationID: myOrg.ID}).
					Return(scalingo.Project{}, forbiddenErr)
			},
			expectedError:      "add project",
			expectedErrorCause: forbiddenErr,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			client := scalingomock.NewMockAPI(gomock.NewController(t))
			test.expectations(ctx, client)

			project, err := createProject(ctx, client, test.params, test.organizationSlug)

			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				if test.expectedErrorCause != nil {
					require.ErrorIs(t, err, test.expectedErrorCause)
				}
				assert.Equal(t, scalingo.Project{}, project)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, project)
		})
	}
}
