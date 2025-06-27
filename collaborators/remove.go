package collaborators

import (
	"context"
	"errors"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	notFound = errors.New("collaborator not found")
)

func Remove(ctx context.Context, app, email string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	collaborator, err := getFromEmail(ctx, client, app, email)
	if err != nil {
		if err == notFound {
			io.Error(email + " is not a collaborator of " + app + ".")
			return nil
		}

		return errgo.Mask(err, errgo.Any)
	}
	err = client.CollaboratorRemove(ctx, app, collaborator.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(email, "has been removed from the collaborators of", app)
	return nil
}

func getFromEmail(ctx context.Context, client *scalingo.Client, app, email string) (scalingo.Collaborator, error) {
	collaborators, err := client.CollaboratorsList(ctx, app)
	if err != nil {
		return scalingo.Collaborator{}, errgo.Mask(err, errgo.Any)
	}
	for _, collaborator := range collaborators {
		if collaborator.Email == email {
			return collaborator, nil
		}
	}
	return scalingo.Collaborator{}, notFound
}
