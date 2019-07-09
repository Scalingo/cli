package collaborators

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
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	collaborator, err := getFromEmail(c, app, email)
	if err != nil {
		if err == notFound {
			io.Error(email + " is not a collaborator of " + app + ".")
			return nil
		} else {
			return errgo.Mask(err, errgo.Any)
		}
	}
	err = c.CollaboratorRemove(app, collaborator.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(email, "has been removed from the collaborators of", app)
	return nil
}

func getFromEmail(c *scalingo.Client, app, email string) (scalingo.Collaborator, error) {
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
