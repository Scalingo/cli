package autoscalers

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v10"
)

func Update(ctx context.Context, app, containerType string, params scalingo.AutoscalerUpdateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	autoscaler, err := getFromContainerType(ctx, c, app, containerType)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			io.Error("Container type " + containerType + " has no autoscaler on the app " + app + ".")
			return nil
		}
		return errors.Wrap(ctx, err, "operation failed")
	}
	_, err = c.AutoscalerUpdate(ctx, app, autoscaler.ID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}

	io.Status("Autoscaler updated on", app, "for", containerType, "containers")
	return nil
}
