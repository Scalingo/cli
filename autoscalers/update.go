package autoscalers

import (
	"context"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors"
)

func Update(ctx context.Context, app, containerType string, params scalingo.AutoscalerUpdateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	autoscaler, err := getFromContainerType(ctx, c, app, containerType)
	if err != nil {
		if errors.ErrgoRoot(err) == ErrNotFound {
			io.Error("Container type " + containerType + " has no autoscaler on the app " + app + ".")
			return nil
		}
		return errgo.Mask(err, errgo.Any)
	}
	_, err = c.AutoscalerUpdate(ctx, app, autoscaler.ID, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler updated on", app, "for", containerType, "containers")
	return nil
}
