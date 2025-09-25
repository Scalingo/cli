package dbng

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	databases, err := c.Preview().DatabasesList(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "list databases")
	}

	io.Warning("This command only displays databases next generation you own.")

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Name", "Type", "Plan", "Role", "Status", "Project"})

	currentUser, err := config.C.CurrentUser(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get current user")
	}

	for _, db := range databases {
		role := utils.AppRole(currentUser, &db.App)

		_ = t.Append([]string{
			db.App.ID,
			db.App.Name,
			db.Addon.AddonProvider.Name,
			db.Addon.Plan.Name,
			string(role),
			string(db.Addon.Status),
			db.App.Project.Name,
		})
	}
	_ = t.Render()

	io.Info("Looking for database addons attached to your applications? Use the 'addons' command")
	io.Info("Example: scalingo --app my-app addons")

	return nil
}
