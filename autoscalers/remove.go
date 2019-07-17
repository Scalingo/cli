package autoscalers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors"
	"gopkg.in/errgo.v1"
)

func Remove(app, containerType string) error {
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

	err = c.AutoscalerRemove(app, autoscaler.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler removed on", app, "for", containerType, "containers")
	return nil
}
