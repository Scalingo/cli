package apps

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-utils/errors/v3"
)

func List(ctx context.Context, renderer renderer.Renderer[[]*scalingo.App], projectSlug string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	apps, err := c.AppsList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list apps")
	}

	filteredApps := filterAppsByProject(apps, projectSlug)
	renderer.SetData(ctx, filteredApps)

	err = renderer.Render(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "render apps list")
	}

	return nil
}

func filterAppsByProject(apps []*scalingo.App, projectSlug string) []*scalingo.App {
	if projectSlug == "" {
		return apps
	}

	filteredApps := make([]*scalingo.App, 0, len(apps))
	for _, app := range apps {
		if app.ProjectSlug() == projectSlug {
			filteredApps = append(filteredApps, app)
		}
	}

	return filteredApps
}
