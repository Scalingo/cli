package autoscalers

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Remove(ctx context.Context, app, containerType string) error {
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

	err = c.AutoscalerRemove(ctx, app, autoscaler.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}

	io.Status("Autoscaler removed on", app, "for", containerType, "containers")
	return nil
}
