package log_drains

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

type addonObject struct {
	Drains    []scalingo.LogDrain
	AddonName string
}

func ListAddon(app string, addonID string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	resources, err := c.AddonsList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list addons")
	}

	addonsToPrint := []addonObject{}
	for _, resource := range resources {
		if addonID == resource.ID || addonID == "" {
			res, err := c.LogDrainsAddonList(app, resource.ID)
			if err != nil {
				return errgo.Notef(err, "fail to list the log drains")
			}
			addonsToPrint = append(addonsToPrint, addonObject{
				AddonName: resource.AddonProvider.Name,
				Drains:    res.Drains,
			})
			if addonID != "" {
				break
			}
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

	drawAddonTable(t, addonsToPrint)
	return nil
}

func drawAddonTable(t *tablewriter.Table, addons []addonObject) {
	addonsLength := len(addons)
	for _, addonsDrains := range addons {
		if len(addonsDrains.Drains) > 1 && addonsLength > 1 {
			t.SetRowLine(true)
		}
		for _, logDrain := range addonsDrains.Drains {
			t.Append([]string{
				addonsDrains.AddonName,
				logDrain.URL,
			})
		}
	}
	t.Render()

}
