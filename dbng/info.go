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

func Info(ctx context.Context, appID string) error {
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

	data := [][]string{
		{"ID", db.App.ID},
		{"Name", db.App.Name},
		{"Type", db.Addon.AddonProvider.Name},
		{"Plan", db.Addon.Plan.Name},
		{"Status", string(db.Database.Status)},
		{"Version", db.Database.ReadableVersion},
		{"Force TLS  ", forceSSL},
		{"Internet Accessibility", internetAccess},
		{"Maintenance window", utils.FormatMaintenanceWindowWithTimezone(db.Database.MaintenanceWindow, time.Local)},
	}
	t := tablewriter.NewWriter(os.Stdout)
	_ = t.Bulk(data)
	_ = t.Render()

	return nil
}
