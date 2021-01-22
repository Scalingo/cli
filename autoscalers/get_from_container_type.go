package autoscalers

import (
	"errors"

	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

var (
	ErrNotFound = errors.New("autoscaler not found")
)

func getFromContainerType(c *scalingo.Client, app, containerType string) (scalingo.Autoscaler, error) {
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
