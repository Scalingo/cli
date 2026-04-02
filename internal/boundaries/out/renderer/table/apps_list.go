package table

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

type appsListRenderer struct {
	currentUser *scalingo.User
	apps        []*scalingo.App
}

func NewAppsList(currentUser *scalingo.User) renderer.Renderer[[]*scalingo.App] {
	return &appsListRenderer{
		currentUser: currentUser,
	}
}

func (r *appsListRenderer) Render(ctx context.Context) error {
	if len(r.apps) == 0 {
		fmt.Println(io.Indent("\nYou haven't created any app yet, create your first application using:\n→ scalingo create <app_name>\n", 2))
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Role", "Status", "Project"})

	for _, app := range r.apps {
		role := utils.AppRole(r.currentUser, app)
		err := t.Append([]string{app.Name, string(role), string(app.Status), app.ProjectSlug()})
		if err != nil {
			return errors.Wrap(ctx, err, "append app to table")
		}
	}

	return t.Render()
}

func (r *appsListRenderer) SetData(ctx context.Context, apps []*scalingo.App) {
	r.apps = apps
}
