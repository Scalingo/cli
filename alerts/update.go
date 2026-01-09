package alerts

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Update(ctx context.Context, app, id string, params scalingo.AlertUpdateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get Scalingo client")
	}

	_, err = c.AlertUpdate(ctx, app, id, params)
	if err != nil {
		return errors.Wrapf(ctx, err, "update alert")
	}

	var msg string
	if params.Disabled != nil {
		if *params.Disabled {
			msg = "Alert disabled"
		} else {
			msg = "Alert enabled"
		}
	} else {
		msg = "Alert updated"
	}
	io.Status(msg)
	return nil
}
