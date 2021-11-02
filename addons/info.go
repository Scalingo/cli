package addons

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/db"
)

func Info(app, addon string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	addonInfo, err := c.AddonShow(app, addon)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	dbInfo, err := db.Show(app,addon)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	forceSsl, internetAccess := "disabled", "disabled"
	if featuresLen := len(dbInfo.Features); featuresLen > 0 {
		if featuresLen == 2 {
			internetAccess = strings.ToLower(dbInfo.Features[1]["status"])
		}
		forceSsl = strings.ToLower(dbInfo.Features[0]["status"])
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