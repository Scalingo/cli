package projects

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Remove(ctx context.Context, projectID string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	err = client.ProjectDelete(ctx, projectID)
	if err != nil {
		return errors.Wrap(ctx, err, "delete project")
	}

	io.Status(projectID, "has been removed")

	return nil
}
