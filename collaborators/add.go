package collaborators

import (
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Add(app, email string) error {
	collaborator, err := api.CollaboratorAdd(app, email)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(collaborator.Email, "has been invited to collaborate to", app)
	return nil
}
