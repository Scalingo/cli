package collaborators

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Add(ctx context.Context, app string, params scalingo.CollaboratorAddParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	collaborator, err := c.CollaboratorAdd(ctx, app, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add collaborator")
	}

	io.Status(collaborator.Email, "has been invited to collaborate to", app)
	return nil
}
