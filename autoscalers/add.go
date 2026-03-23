package autoscalers

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Add(ctx context.Context, app string, params scalingo.AutoscalerAddParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	autoscaler, err := c.AutoscalerAdd(ctx, app, params)
	if err != nil {
		return errors.Wrapf(ctx, err, "create autoscaler on app %s", app)
	}

	io.Status("Autoscaler created on", app, "for", autoscaler.ContainerType, "containers")
	return nil
}
