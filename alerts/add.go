package alerts

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v5"
)

func Add(ctx context.Context, app string, params scalingo.AlertAddParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	a, err := c.AlertAdd(ctx, app, params)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Alert created for the container type", a.ContainerType)
	return nil
}
