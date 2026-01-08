package collaborators

import (
	"context"
	stderrors "errors"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

var (
	errNotFound = stderrors.New("collaborator not found")
)

func Remove(ctx context.Context, app, email string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	collaborator, err := getFromEmail(ctx, client, app, email)
	if errors.Is(err, errNotFound) {
		io.Error(email + " is not a collaborator of " + app + ".")
		return nil
	} else if err != nil {
		return errors.Wrap(ctx, err, "get from email")
	}

	err = client.CollaboratorRemove(ctx, app, collaborator.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "remove collaborator")
	}

	io.Status(email, "has been removed from the collaborators of", app)
	return nil
}

func getFromEmail(ctx context.Context, client *scalingo.Client, app, email string) (scalingo.Collaborator, error) {
	collaborators, err := client.CollaboratorsList(ctx, app)
	if err != nil {
		return scalingo.Collaborator{}, errors.Wrap(ctx, err, "list collaborators")
	}
	for _, collaborator := range collaborators {
		if collaborator.Email == email {
			return collaborator, nil
		}
	}
	return scalingo.Collaborator{}, errNotFound
}
