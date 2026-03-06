package table

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-utils/errors/v3"
)

type appsListRenderer struct {
	currentUser *scalingo.User
}

func NewAppsList(currentUser *scalingo.User) apps.ListRenderer {
	return appsListRenderer{
		currentUser: currentUser,
	}
}

func (r appsListRenderer) Render(ctx context.Context, apps []*scalingo.App) error {
	if len(apps) == 0 {
		fmt.Println(io.Indent("\nYou haven't created any app yet, create your first application using:\n→ scalingo create <app_name>\n", 2))
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Role", "Status", "Project"})

	for _, app := range apps {
		role := utils.AppRole(r.currentUser, app)
		err := t.Append([]string{app.Name, string(role), string(app.Status), app.ProjectSlug()})
		if err != nil {
			return errors.Wrap(ctx, err, "append app to table")
		}
	}

	return t.Render()
}
