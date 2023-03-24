package autoscalers

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	errors "github.com/Scalingo/go-utils/errors/v2"
)

func Remove(ctx context.Context, app, containerType string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	autoscaler, err := getFromContainerType(ctx, c, app, containerType)
	if err != nil {
		if errors.RootCause(err) == ErrNotFound {
			io.Error("Container type " + containerType + " has no autoscaler on the app " + app + ".")
			return nil
		}
		return errgo.Mask(err, errgo.Any)
	}

	err = c.AutoscalerRemove(ctx, app, autoscaler.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler removed on", app, "for", containerType, "containers")
	return nil
}
