package collaborators

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

const (
	CollaboratorOwner = "owner"
)

func List(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	collaborators, err := c.CollaboratorsList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list collaborators")
	}

	scapp, err := c.AppsShow(app)
	if err != nil {
		return errgo.Notef(err, "fail to get application information")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Email", "Username", "Status"})

	t.Append([]string{scapp.Owner.Email, scapp.Owner.Username, CollaboratorOwner})
	for _, collaborator := range collaborators {
		t.Append([]string{collaborator.Email, collaborator.Username, string(collaborator.Status)})
	}
	t.Render()
	return nil
}
