package collaborators

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

const (
	CollaboratorOwner = "owner"
)

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	collaborators, err := c.CollaboratorsList(ctx, app)
	if err != nil {
		return errgo.Notef(err, "fail to list collaborators")
	}

	scapp, err := c.AppsShow(ctx, app)
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
