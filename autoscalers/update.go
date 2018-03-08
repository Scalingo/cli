package autoscalers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	errgo "gopkg.in/errgo.v1"
)

func Update(app, containerType string, params scalingo.AutoscalerUpdateParams) error {
	autoscaler, err := getFromContainerType(app, containerType)
	if err != nil {
		if err == errNotFound {
			io.Error("Container type " + containerType + " has no autoscaler on the app " + app + ".")
			return nil
		}
		return errgo.Mask(err, errgo.Any)
	}

	c := config.ScalingoClient()
	_, err = c.AutoscalerUpdate(app, autoscaler.ID, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler updated on", app, "for", containerType, "containers")
	return nil
}
