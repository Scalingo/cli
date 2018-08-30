package autoscalers

import (
	"errors"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

var (
	ErrNotFound = errors.New("autoscaler not found")
)

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
	return scalingo.Autoscaler{}, ErrNotFound
}
