package alerts

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v7"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Add(ctx context.Context, app string, params scalingo.AlertAddParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	a, err := c.AlertAdd(ctx, app, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add alert")
	}

	io.Status("Alert created for the container type", a.ContainerType)
	return nil
}
