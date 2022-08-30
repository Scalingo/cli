package alerts

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Remove(ctx context.Context, app, id string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	err = c.AlertRemove(ctx, app, id)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The alert has been deleted")
	return nil
}
