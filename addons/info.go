package addons

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/db"
)

func Info(ctx context.Context, app, addon string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	addonInfo, err := c.AddonShow(ctx, app, addon)
	if err != nil {
		return errgo.Notef(err, "fail to get addon information")
	}

	dbInfo, err := db.Show(ctx, app, addon)
	if err != nil {
		return errgo.Notef(err, "fail to get database information")
	}

	forceSsl, internetAccess := "disabled", "disabled"
	for i := range dbInfo.Features {
		if dbInfo.Features[i]["name"] == "force-ssl" {
			forceSsl = strings.ToLower(dbInfo.Features[i]["status"])
		} else if dbInfo.Features[i]["name"] == "publicly-available" {
			internetAccess = strings.ToLower(dbInfo.Features[i]["status"])
		}
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Append([]string{"Database Type", fmt.Sprintf("%v", dbInfo.TypeName)})
	t.Append([]string{"Version", fmt.Sprintf("%v", dbInfo.ReadableVersion)})
	t.Append([]string{"Status", fmt.Sprintf("%v", addonInfo.Status)})
	t.Append([]string{"Plan", fmt.Sprintf("%v", addonInfo.Plan.Name)})
	t.Append([]string{"Force TLS", forceSsl})
	t.Append([]string{"Internet Accessibility", internetAccess})

	t.Render()

	return nil
}
