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

	databases, err := listWithOrganizations(ctx, c, c.Preview())
	if err != nil {
		return errors.Wrap(ctx, err, "list databases with organizations")
	}

	io.Warning("This command only displays databases with dedicated resources (DR).")

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

// listWithOrganizations returns the personal databases followed by the databases of each organization of the user.
func listWithOrganizations(ctx context.Context, organizationsService scalingo.OrganizationsService, databasesService scalingo.DatabasesPreviewService) ([]listedDatabase, error) {
	personalDatabases, err := databasesService.DatabasesList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list databases")
	}

	organizations, err := organizationsService.OrganizationsList(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list organizations")
	}

	databases := make([]listedDatabase, 0, len(personalDatabases))
	for _, database := range personalDatabases {
		databases = append(databases, listedDatabase{database: database})
	}

	for _, organization := range organizations {
		organizationDatabases, err := databasesService.DatabasesListByOrganization(ctx, organization.ID)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "list organization databases")
		}
		for _, database := range organizationDatabases {
			databases = append(databases, listedDatabase{database: database, organizationSlug: organization.Slug})
		}
	}

	return databases, nil
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
