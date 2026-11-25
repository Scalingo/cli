package apps

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

// ListedApp is an app with the organization in which it has been listed.
// OrganizationID and OrganizationSlug are empty for an app of the personal scope.
type ListedApp struct {
	App              *scalingo.App
	OrganizationID   string
	OrganizationSlug string
}

func List(ctx context.Context, renderer renderer.Renderer[[]ListedApp], projectSlug string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	apps, err := listWithOrganizations(ctx, c)
	if err != nil {
		return errors.Wrap(ctx, err, "list apps with organizations")
	}

	filteredApps := filterAppsByProject(apps, projectSlug)
	renderer.SetData(ctx, filteredApps)

	err = renderer.Render(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "render apps list")
	}

	return nil
}

// listWithOrganizations returns the personal apps followed by the apps of each organization of the user.
func listWithOrganizations(ctx context.Context, client scalingo.API) ([]ListedApp, error) {
	personalApps, err := client.AppsList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list apps")
	}

	organizations, err := client.OrganizationsList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list organizations")
	}

	// A nil personal list stays nil when no organization app is added, to keep the historical {"apps":null} JSON output.
	var apps []ListedApp
	if personalApps != nil {
		apps = make([]ListedApp, 0, len(personalApps))
	}
	for _, app := range personalApps {
		apps = append(apps, ListedApp{App: app})
	}

	for _, organization := range organizations {
		organizationApps, err := client.AppsListByOrganization(ctx, organization.ID)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "list organization apps")
		}
		for _, app := range organizationApps {
			apps = append(apps, ListedApp{App: app, OrganizationID: organization.ID, OrganizationSlug: organization.Slug})
		}
	}

	return apps, nil
}

// filterAppsByProject keeps the apps whose project slug, under the format <ownerUsername>/<projectName>, matches projectSlug.
func filterAppsByProject(apps []ListedApp, projectSlug string) []ListedApp {
	if projectSlug == "" {
		return apps
	}

	filteredApps := make([]ListedApp, 0, len(apps))
	for _, listedApp := range apps {
		if listedApp.App.ProjectSlug() == projectSlug {
			filteredApps = append(filteredApps, listedApp)
		}
	}

	return filteredApps
}
