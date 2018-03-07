package autoscalers

import (
	"errors"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
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

	c := config.ScalingoClient()
	err = c.CollaboratorRemove(app, collaborator.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(email, "has been removed from the collaborators of", app)
	return nil
}

func getFromEmail(app, email string) (scalingo.Collaborator, error) {
	c := config.ScalingoClient()
	collaborators, err := c.CollaboratorsList(app)
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
