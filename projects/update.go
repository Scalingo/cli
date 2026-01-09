package projects

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Update(ctx context.Context, projectID string, params scalingo.ProjectUpdateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	project, err := c.ProjectUpdate(ctx, projectID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "update project")
	}

	io.Status(project.Name, "has been updated")

	return nil
}
