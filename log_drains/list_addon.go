package log_drains

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func ListAddon(app string, addonID string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	var addons []scalingo.LogDrainsAddonRes

	if addonID != "" {
		addonsDrains, err := c.LogDrainsAddonList(app, addonID)
		if err != nil {
			return errgo.Notef(err, "fail to list the log drains")
		}
		addons = make([]scalingo.LogDrainsAddonRes, 1)
		addons[0] = addonsDrains

	} else {
		resources, err := c.AddonsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list addons")
		}

		addons = make([]scalingo.LogDrainsAddonRes, len(resources))
		for index, resource := range resources {
			addonsDrains, err := c.LogDrainsAddonList(app, resource.ID)
			if err != nil {
				return errgo.Notef(err, "fail to list the log drains")
			}

			addons[index] = addonsDrains
		}
	}

	t := tablewriter.NewWriter(os.Stdout)

	t.SetHeader([]string{"Name", "URL"})
	t.SetAutoMergeCells(true)

	if addonID == "" {
		appDrains, err := c.LogDrainsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list the log drains")
		}

		for _, logDrain := range appDrains {
			t.Append([]string{
				app,
				logDrain.URL,
			})
		}
	}

	drawAddonTable(t, addons)
	return nil
}

func drawAddonTable(t *tablewriter.Table, addons []scalingo.LogDrainsAddonRes) {

	addonsLength := len(addons)
	for _, addonsDrains := range addons {
		if len(addonsDrains.Drains) > 1 && addonsLength > 1 {
			t.SetRowLine(true)
		}
		for _, logDrain := range addonsDrains.Drains {
			t.Append([]string{
				addonsDrains.Addon.Name,
				logDrain.URL,
			})
		}
	}
	t.Render()

}
