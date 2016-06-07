package collaborators

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func List(app string) error {
	c := config.ScalingoClient()
	collaborators, err := c.CollaboratorsList(app)
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
