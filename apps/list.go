package apps

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

const (
	roleOwner        = "owner"
	roleCollaborator = "collaborator"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	apps, err := c.AppsList(ctx)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if len(apps) == 0 {
		fmt.Println(io.Indent("\nYou haven't created any app yet, create your first application using:\n→ scalingo create <app_name>\n", 2))
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Role", "Status", "Project"})

	currentUser, err := config.C.CurrentUser(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get current user")
	}

	for _, app := range apps {
		role := roleCollaborator
		if app.Owner.Email == currentUser.Email {
			role = roleOwner
		}

		t.Append([]string{app.Name, role, string(app.Status), app.Project.Name})
	}
	t.Render()

	return nil
}
