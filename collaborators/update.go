package collaborators

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
)

func Update(ctx context.Context, app, collaboratorID string, params scalingo.CollaboratorUpdateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	collaborator, err := c.CollaboratorUpdate(ctx, app, collaboratorID, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(collaborator.Email, "role has been updated for", app)
	return nil
}
