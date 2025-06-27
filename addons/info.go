// Package addons gather all the command handlers related to addons in general
// Some specific treatment are done for database addons (like displaying
// database feature status)
package addons

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/db"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

// Info is the command handler displaying static information about one given addon
func Info(ctx context.Context, app, addon string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	addonInfo, err := c.AddonShow(ctx, app, addon)
	if err != nil {
		return errors.Wrap(ctx, err, "get addon information")
	}

	isDatabase := slices.ContainsFunc(db.SupportedDatabases, func(s string) bool {
		return strings.EqualFold(s, addonInfo.AddonProvider.ID)
	})

	dbInfo := [][]string{}
	if isDatabase {
		dbInfo, err = getDatabaseInfo(ctx, c, app, addon)
		if err != nil {
			return errors.Wrap(ctx, err, "get database information")
		}
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Append([]string{"Addon Provider", addonInfo.AddonProvider.Name})
	t.Append([]string{"Plan", addonInfo.Plan.Name})
	t.Append([]string{"Status", fmt.Sprintf("%v", addonInfo.Status)})
	for _, line := range dbInfo {
		t.Append(line)
	}

	t.Render()

	return nil
}

func getDatabaseInfo(ctx context.Context, c *scalingo.Client, app, addon string) ([][]string, error) {
	dbInfo, err := c.DatabaseShow(ctx, app, addon)
	if err != nil {
		return [][]string{}, errors.Wrap(ctx, err, "get database information")
	}

	forceSsl, internetAccess := "disabled", "disabled"
	for _, feature := range dbInfo.Features {
		if feature.Name == "force-ssl" {
			forceSsl = strings.ToLower(string(feature.Status))
		} else if feature.Name == "publicly-available" {
			internetAccess = strings.ToLower(string(feature.Status))
		}
	}

	lines := [][]string{
		{"Database Type", dbInfo.TypeName},
		{"Version", dbInfo.ReadableVersion},
		{"Force TLS", forceSsl},
		{"Internet Accessibility", internetAccess},
		{"Maintenance window", utils.FormatMaintenanceWindowWithTimezone(dbInfo.MaintenanceWindow, time.Local)},
	}
	return lines, nil
}
