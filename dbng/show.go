package dbng

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Show(ctx context.Context, appID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	db, err := c.Preview().DatabaseShow(ctx, appID)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	forceSSL, internetAccess := "disabled", "disabled"
	for _, feature := range db.Database.Features {
		switch feature.Name {
		case "force-ssl":
			forceSSL = strings.ToLower(string(feature.Status))
		case "publicly-available":
			internetAccess = strings.ToLower(string(feature.Status))
		}
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Append("ID", db.App.ID)
	t.Append("Name", db.App.Name)
	t.Append("Type", db.Addon.AddonProvider.Name)
	t.Append("Plan", db.Addon.Plan.Name)
	t.Append("Status", db.Database.Status)
	t.Append("Version", db.Database.ReadableVersion)
	t.Append("Force TLS  ", forceSSL)
	t.Append("Internet Accessibility", internetAccess)
	t.Append("Maintenance window", utils.FormatMaintenanceWindowWithTimezone(db.Database.MaintenanceWindow, time.Local))
	t.Render()

	return nil
}
