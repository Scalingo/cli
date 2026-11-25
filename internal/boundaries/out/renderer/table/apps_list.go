package table

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/internal/boundaries/out/renderer"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

type appsListRenderer struct {
	currentUser *scalingo.User
	listedApps  []apps.ListedApp
}

func NewAppsList(currentUser *scalingo.User) renderer.Renderer[[]apps.ListedApp] {
	return &appsListRenderer{
		currentUser: currentUser,
	}
}

func (r *appsListRenderer) Render(ctx context.Context) error {
	if len(r.listedApps) == 0 {
		fmt.Println(io.Indent("\nYou haven't created any app yet in this project, create your first application using:\n→ scalingo create <app_name>\n", 2))
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Role", "Status", "Project"})

	for _, listed := range r.listedApps {
		app := listed.App
		role := utils.AppRole(r.currentUser, app)
		err := t.Append([]string{app.Name, string(role), string(app.Status), appProject(listed)})
		if err != nil {
			return errors.Wrap(ctx, err, "append app to table")
		}
	}

	return t.Render()
}

func (r *appsListRenderer) SetData(ctx context.Context, listedApps []apps.ListedApp) {
	r.listedApps = listedApps
}

// appProject returns <ownerUsername>/<projectName> for a personal app and <organizationSlug>/<projectName> for an app of an organization.
func appProject(listed apps.ListedApp) string {
	if listed.Organization == nil {
		return listed.App.ProjectSlug()
	}
	return utils.OrganizationProjectLabel(listed.Organization.Slug, listed.App.Project.Name)
}
