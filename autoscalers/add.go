package autoscalers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Add(app, email string) error {
	c := config.ScalingoClient()
	collaborator, err := c.CollaboratorAdd(app, email)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status(collaborator.Email, "has been invited to collaborate to", app)
	return nil
}
