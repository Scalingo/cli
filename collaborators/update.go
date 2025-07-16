package collaborators

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Update(ctx context.Context, app, collaboratorID string, params scalingo.CollaboratorUpdateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to get Scalingo client")
	}
	collaborator, err := c.CollaboratorUpdate(ctx, app, collaboratorID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "update collaborator")
	}

	io.Status(collaborator.Email, "role has been updated for", app)
	return nil
}
