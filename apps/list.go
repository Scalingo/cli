package apps

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

const (
	roleOwner        = "owner"
	roleCollaborator = "collaborator"
)

func List(ctx context.Context, projectSlug string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	apps, err := c.AppsList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list apps")
	}

	if len(apps) == 0 {
		fmt.Println(io.Indent("\nYou haven't created any app yet, create your first application using:\n→ scalingo create <app_name>\n", 2))
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Role", "Status", "Project"})

	currentUser, err := config.C.CurrentUser(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to get current user")
	}

	for _, app := range apps {
		// If a filter was set but the app is not in the project, skip to the next one.
		if projectSlug != "" && projectSlug != app.ProjectSlug() {
			continue
		}

		role := roleCollaborator
		if app.Owner.Email == currentUser.Email {
			role = roleOwner
		}

		_ = t.Append([]string{app.Name, role, string(app.Status), app.ProjectSlug()})
	}
	_ = t.Render()

	return nil
}
