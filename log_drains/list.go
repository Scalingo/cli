package log_drains

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

type printableDrains struct {
	DrainURLs []scalingo.LogDrain
	AppName   string
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

	appToPrint := []printableDrains{}

	if opts.AddonID == "" {
		logDrains, err := c.LogDrainsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list the log drains")
		}
		if len(logDrains) > 0 {
			appToPrint = append(appToPrint, printableDrains{
				AppName:   app,
				DrainURLs: logDrains,
			})
		}
	}
	if opts.AddonID != "" || opts.WithAddons {
		addons, err := c.AddonsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list addons")
		}

		for _, addon := range addons {
			if opts.AddonID == addon.ID || opts.WithAddons {
				res, err := c.LogDrainsAddonList(app, addon.ID)
				if err != nil {
					return errgo.Notef(err, "fail to list the log drains of an addon")
				}
				if len(res.Drains) > 0 {
					appToPrint = append(appToPrint, printableDrains{
						AppName:   addon.AddonProvider.Name,
						DrainURLs: res.Drains,
					})
				}

				if !opts.WithAddons {
					break
				}
			}
		}
	}

	drawDrainsTable(appToPrint)
	return nil
}

func drawDrainsTable(drains []printableDrains) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "URL"})
	t.SetAutoMergeCells(true)

	objLength := len(drains)

	for _, printableDrain := range drains {
		if len(printableDrain.DrainURLs) > 1 && objLength > 1 {
			t.SetRowLine(true)
		}
		for _, drain := range printableDrain.DrainURLs {
			t.Append([]string{
				printableDrain.AppName,
				drain.URL,
			})
		}
	}
	t.Render()
}
