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

	personalApps, err := c.AppsList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list apps")
	}

	apps := make([]ListedApp, 0, len(personalApps))
	for _, app := range personalApps {
		apps = append(apps, ListedApp{App: app})
	}

	filteredApps := filterAppsByProject(apps, projectSlug)
	renderer.SetData(ctx, filteredApps)

	err = renderer.Render(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "render apps list")
	}

	return nil
}

// filterAppsByProject keeps the apps whose project slug, under the format <ownerUsername>/<projectName>, matches projectSlug.
func filterAppsByProject(apps []ListedApp, projectSlug string) []ListedApp {
	if projectSlug == "" {
		return apps
	}

	filteredApps := make([]ListedApp, 0, len(apps))
	for _, listed := range apps {
		if listed.App.ProjectSlug() == projectSlug {
			filteredApps = append(filteredApps, listed)
		}
	}

	return filteredApps
}
