package orgs

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func List(ctx context.Context, r renderer.Renderer[[]scalingo.Organization]) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	return list(ctx, c, r)
}

func list(ctx context.Context, service scalingo.OrganizationsService, r renderer.Renderer[[]scalingo.Organization]) error {
	organizations, err := service.OrganizationsList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list organizations")
	}

	r.SetData(ctx, organizations)

	err = r.Render(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "render organizations list")
	}

	return nil
}
