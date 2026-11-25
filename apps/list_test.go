package apps

import (
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

	// The personal scope and the organizations contain projects with the same name.
	const projectName = "app"
	owner := scalingo.Owner{Username: "john"}
	personalWeb := &scalingo.App{Name: "web", Owner: owner}
	personalWeb.Project.Name = projectName
	acmeAPI := &scalingo.App{Name: "api", Owner: owner}
	acmeAPI.Project.Name = projectName
	acmeZoo := &scalingo.App{Name: "zoo", Owner: owner}
	acmeZoo.Project.Name = "back"
	globexAdmin := &scalingo.App{Name: "admin", Owner: owner}
	globexAdmin.Project.Name = projectName

	tests := map[string]struct {
		expectations  func(ctx context.Context, client *scalingomock.MockAPI)
		expected      []ListedApp
		expectedError string
		expectedCause error
	}{
		"it lists only the personal apps of a user without organization": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return([]*scalingo.App{personalWeb}, nil)
				client.EXPECT().OrganizationsList(ctx).Return(nil, nil)
			},
			expected: []ListedApp{{App: personalWeb}},
		},
		"it lists the personal apps, then the apps of each organization in the order returned": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return([]*scalingo.App{personalWeb}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				client.EXPECT().AppsListByOrganization(ctx, acme.ID).Return([]*scalingo.App{acmeAPI, acmeZoo}, nil)
				client.EXPECT().AppsListByOrganization(ctx, globex.ID).Return([]*scalingo.App{globexAdmin}, nil)
			},
			expected: []ListedApp{
				{App: personalWeb},
				{App: acmeAPI, Organization: &acme},
				{App: acmeZoo, Organization: &acme},
				{App: globexAdmin, Organization: &globex},
			},
		},
		"it continues with the other organizations when one has no app": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return(nil, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				client.EXPECT().AppsListByOrganization(ctx, acme.ID).Return([]*scalingo.App{}, nil)
				client.EXPECT().AppsListByOrganization(ctx, globex.ID).Return([]*scalingo.App{globexAdmin}, nil)
			},
			expected: []ListedApp{{App: globexAdmin, Organization: &globex}},
		},
		"it keeps a nil list when no scope returns any app": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return(nil, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme}, nil)
				client.EXPECT().AppsListByOrganization(ctx, acme.ID).Return([]*scalingo.App{}, nil)
			},
			expected: nil,
		},
		"it keeps an empty list when the personal list is empty": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return([]*scalingo.App{}, nil)
				client.EXPECT().OrganizationsList(ctx).Return(nil, nil)
			},
			expected: []ListedApp{},
		},
		"it returns an error when the personal apps cannot be listed": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return(nil, stderrors.New("boom"))
			},
			expectedError: "list apps",
		},
		"it returns an error when the organizations cannot be listed": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return([]*scalingo.App{personalWeb}, nil)
				client.EXPECT().OrganizationsList(ctx).Return(nil, stderrors.New("boom"))
			},
			expectedError: "list organizations",
		},
		"it returns an error when organization apps cannot be listed": {
			expectations: func(ctx context.Context, client *scalingomock.MockAPI) {
				client.EXPECT().AppsList(ctx).Return([]*scalingo.App{personalWeb}, nil)
				client.EXPECT().OrganizationsList(ctx).Return([]scalingo.Organization{acme, globex}, nil)
				client.EXPECT().AppsListByOrganization(ctx, acme.ID).Return([]*scalingo.App{acmeAPI}, nil)
				client.EXPECT().AppsListByOrganization(ctx, globex.ID).Return(nil, errForbidden)
			},
			expectedError: "list organization apps",
			expectedCause: errForbidden,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			client := scalingomock.NewMockAPI(gomock.NewController(t))
			test.expectations(ctx, client)

			apps, err := listWithOrganizations(ctx, client)

			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				if test.expectedCause != nil {
					require.ErrorIs(t, err, test.expectedCause)
				}
				assert.Nil(t, apps)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, apps)
			}
		})
	}
}

func TestFilterAppsByProject(t *testing.T) {
	const projectName = "app"
	organization := &scalingo.Organization{ID: "org-acme", Slug: "acme"}
	owner := scalingo.Owner{Username: "john"}
	personalApp := &scalingo.App{Name: "personal-web", Owner: owner}
	personalApp.Project.Name = projectName
	organizationApp := &scalingo.App{Name: "acme-web", Owner: owner}
	organizationApp.Project.Name = projectName
	otherProjectApp := &scalingo.App{Name: "acme-worker", Owner: owner}
	otherProjectApp.Project.Name = "back"

	apps := []ListedApp{
		{App: personalApp},
		{App: organizationApp, Organization: organization},
		{App: otherProjectApp, Organization: organization},
	}

	tests := map[string]struct {
		projectSlug string
		expected    []ListedApp
	}{
		"it keeps all the apps without filter": {
			expected: apps,
		},
		"it filters the apps of all the scopes with the <ownerUsername>/<projectName> format": {
			projectSlug: owner.Username + "/" + projectName,
			expected:    []ListedApp{apps[0], apps[1]},
		},
		"it returns an empty list when the filter matches no app": {
			projectSlug: owner.Username + "/unknown",
			expected:    []ListedApp{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, filterAppsByProject(apps, test.projectSlug))
		})
	}
}
