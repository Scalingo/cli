package dbdr

import (
	"context"
	stdio "io"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

// listedDatabase is a database with the slug of the organization in which it has been listed, empty for the personal scope.
type listedDatabase struct {
	database         scalingo.DatabaseNG
	organizationSlug string
}

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	personalDatabases, err := c.Preview().DatabasesList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list databases")
	}

	databases := make([]listedDatabase, 0, len(personalDatabases))
	for _, database := range personalDatabases {
		databases = append(databases, listedDatabase{database: database})
	}

	io.Warning("This command only displays the databases with dedicated resources (DR) you own or collaborate on.")

	currentUser, err := config.C.CurrentUser(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get current user")
	}

	err = render(os.Stdout, currentUser, databases)
	if err != nil {
		return errors.Wrap(ctx, err, "render listed databases")
	}

	io.Info("Looking for database addons attached to your applications? Use the 'addons' command")
	io.Info("Example: scalingo --app my-app addons")

	return nil
}

func render(w stdio.Writer, currentUser *scalingo.User, databases []listedDatabase) error {
	t := tablewriter.NewWriter(w)
	t.Header([]string{"ID", "Name", "Type", "Plan", "Role", "Status", "Project"})

	for _, listed := range databases {
		db := listed.database
		role := utils.AppRole(currentUser, &db.App)

		_ = t.Append([]string{
			db.ID,
			db.Name,
			db.Technology,
			db.Plan,
			string(role),
			string(db.Database.Status),
			utils.OrganizationProjectLabel(listed.organizationSlug, db.App.Project.Name),
		})
	}

	return t.Render()
}
