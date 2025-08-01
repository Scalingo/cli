package collaborators

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Update(ctx context.Context, app, email string, params scalingo.CollaboratorUpdateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	collaborator, err := getFromEmail(ctx, c, app, email)
	if errors.Is(err, errNotFound) {
		io.Error(email + " is not a collaborator of " + app + ".")
		return nil
	} else if err != nil {
		return errors.Wrap(ctx, err, "get from email")
	}

	collaborator, err = c.CollaboratorUpdate(ctx, app, collaborator.ID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "update collaborator")
	}

	io.Status(collaborator.Email, "role has been updated for", app)
	return nil
}
