package autoscalers

import (
	"errors"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

var (
	errNotFound = errors.New("autoscaler not found")
)

func Remove(app, containerType string) error {
	autoscaler, err := getFromContainerType(app, containerType)
	if err != nil {
		if err == errNotFound {
			io.Error("Container type " + containerType + " has no autoscaler on the app " + app + ".")
			return nil
		}
		return errgo.Mask(err, errgo.Any)
	}

	c := config.ScalingoClient()
	err = c.AutoscalerRemove(app, autoscaler.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Autoscaler removed on", app, "for", containerType, "containers")
	return nil
}

func getFromContainerType(app, containerType string) (scalingo.Autoscaler, error) {
	c := config.ScalingoClient()
	autoscalers, err := c.AutoscalersList(app)
	if err != nil {
		return scalingo.Autoscaler{}, errgo.Mask(err, errgo.Any)
	}

	for _, autoscaler := range autoscalers {
		if autoscaler.ContainerType == containerType {
			return autoscaler, nil
		}
	}
	return scalingo.Autoscaler{}, errNotFound
}
