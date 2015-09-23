package collaborators

import (
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/go-scalingo"
)

func List(app string) error {
	collaborators, err := scalingo.CollaboratorsList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Email", "Username", "Status"})

	for _, collaborator := range collaborators {
		t.Append([]string{collaborator.Email, collaborator.Username, collaborator.Status})
	}
	t.Render()
	return nil
}
