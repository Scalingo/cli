package collaborators

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Add(ctx context.Context, app, email string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	collaborator, err := c.CollaboratorAdd(ctx, app, email)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(collaborator.Email, "has been invited to collaborate to", app)
	return nil
}
