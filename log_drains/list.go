package log_drains

import (
	"os"
	"strings"

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

	logDrains, err := c.LogDrainsList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the log drains")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "URL"})
	t.SetAutoMergeCells(true)

	if opts.AddonID == "" && !opts.WithAddons {
		for _, logDrain := range logDrains {
			t.Append([]string{
				app,
				logDrain.URL,
			})
		}
	} else {
		addons, err := c.AddonsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list addons")
		}

		addonsToPrint := []addonObject{}
		for _, addon := range addons {
			if opts.AddonID == addon.ID || opts.WithAddons {
				res, err := c.LogDrainsAddonList(app, addon.ID)
				if err != nil {
					return errgo.Notef(err, "fail to list the log drains")
				}
				addonsToPrint = append(addonsToPrint, addonObject{
					AddonName: addon.AddonProvider.Name,
					Drains:    res.Drains,
				})
				if !opts.WithAddons {
					break
				}
			}
		}
		drawAddonTable(t, addonsToPrint, opts.WithAddons)
	}

	t.Render()
	return nil
}

func drawAddonTable(t *tablewriter.Table, addons []addonObject, drawSeparateLines bool) {
	var longestName int
	var longestURL int
	for _, addonsDrains := range addons {
		nameLen := len(addonsDrains.AddonName)
		if longestName < nameLen {
			longestName = nameLen
		}
		for _, logDrain := range addonsDrains.Drains {
			URLLen := len(logDrain.URL)
			if longestURL < URLLen {
				longestURL = URLLen
			}
		}
	}
	for index, addonsDrains := range addons {
		if drawSeparateLines && index != 0 {
			t.Append([]string{
				strings.Repeat("-", longestName),
				strings.Repeat("-", longestURL),
			})
		}
		for _, logDrain := range addonsDrains.Drains {
			t.Append([]string{
				addonsDrains.AddonName,
				logDrain.URL,
			})
		}
	}
}
