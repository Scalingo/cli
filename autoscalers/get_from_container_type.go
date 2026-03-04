package autoscalers

import (
	"context"
	"errors"

	"github.com/Scalingo/go-scalingo/v10"
)

var (
	ErrNotFound = errors.New("autoscaler not found")
)

func getFromContainerType(ctx context.Context, c *scalingo.Client, app, containerType string) (scalingo.Autoscaler, error) {
	autoscalers, err := c.AutoscalersList(ctx, app)
	if err != nil {
		return scalingo.Autoscaler{}, err
	}

	for _, autoscaler := range autoscalers {
		if autoscaler.ContainerType == containerType {
			return autoscaler, nil
		}
	}
	return scalingo.Autoscaler{}, ErrNotFound
}
