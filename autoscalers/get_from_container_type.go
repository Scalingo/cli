package autoscalers

import (
	"context"
	stderrors "errors"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

var (
	ErrNotFound = stderrors.New("autoscaler not found")
)

func getFromContainerType(ctx context.Context, c *scalingo.Client, app, containerType string) (scalingo.Autoscaler, error) {
	autoscalers, err := c.AutoscalersList(ctx, app)
	if err != nil {
		return scalingo.Autoscaler{}, errors.Wrap(ctx, err, "list autoscalers")
	}

	for _, autoscaler := range autoscalers {
		if autoscaler.ContainerType == containerType {
			return autoscaler, nil
		}
	}
	return scalingo.Autoscaler{}, ErrNotFound
}
