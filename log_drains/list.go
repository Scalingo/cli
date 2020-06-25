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

type ListAddonOpts struct {
	WithAddons bool
	AddonID    string
}

func List(app string, opts ListAddonOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "URL"})
	t.SetAutoMergeCells(true)

	if opts.AddonID == "" {
		logDrains, err := c.LogDrainsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list the log drains")
		}

		for _, logDrain := range logDrains {
			t.Append([]string{
				app,
				logDrain.URL,
			})
		}
	}
	if opts.AddonID != "" || opts.WithAddons {
		addons, err := c.AddonsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list addons")
		}

		addonsToPrint := []addonObject{}
		for _, addon := range addons {
			if opts.AddonID == addon.ID || opts.WithAddons {
				res, err := c.LogDrainsAddonList(app, addon.ID)
				if err != nil {
					return errgo.Notef(err, "fail to list the log drains of an addon")
				}
				if len(res.Drains) > 0 {
					addonsToPrint = append(addonsToPrint, addonObject{
						AddonName: addon.AddonProvider.Name,
						Drains:    res.Drains,
					})
				}

				if !opts.WithAddons {
					break
				}
			}
		}
		drawAddonTable(t, addonsToPrint)
	}

	t.Render()
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
}
