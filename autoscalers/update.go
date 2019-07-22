package autoscalers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-utils/errors"
	errgo "gopkg.in/errgo.v1"
)

func Update(app, containerType string, params scalingo.AutoscalerUpdateParams) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	autoscaler, err := getFromContainerType(c, app, containerType)
	if err != nil {
		if errors.ErrgoRoot(err) == ErrNotFound {
			io.Error("Container type " + containerType + " has no autoscaler on the app " + app + ".")
			return nil
		}
		return errgo.Mask(err, errgo.Any)
	}
	_, err = c.AutoscalerUpdate(app, autoscaler.ID, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler updated on", app, "for", containerType, "containers")
	return nil
}
