package apps

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func ProjectSet(ctx context.Context, appName, projectID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	_, err = c.AppsSetProject(ctx, appName, projectID)
	if err != nil {
		return errors.Wrap(ctx, err, "set project for app")
	}

	io.Statusf("Project ID has been set to %s on %s\n", projectID, appName)

	return nil
}
