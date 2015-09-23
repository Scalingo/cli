package collaborators

import (
	"errors"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/io"
)

var (
	notFound = errors.New("collaborator not found")
)

func Remove(app, email string) error {
	collaborator, err := getFromEmail(app, email)
	if err != nil {
		if err == notFound {
			io.Error(email + " is not a collaborator of " + app + ".")
			return nil
		} else {
			return errgo.Mask(err, errgo.Any)
		}
	}

	err = scalingo.CollaboratorRemove(app, collaborator.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(email, "has been removed from the collaborators of", app)
	return nil
}

func getFromEmail(app, email string) (scalingo.Collaborator, error) {
	collaborators, err := scalingo.CollaboratorsList(app)
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
