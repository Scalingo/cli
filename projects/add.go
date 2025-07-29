package projects

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Add(ctx context.Context, params scalingo.ProjectAddParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	project, err := c.ProjectAdd(ctx, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add project")
	}
	io.Status(project.Name, "has been created")

	return nil
}
