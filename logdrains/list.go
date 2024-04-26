package logdrains

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v7"
)

type printableDrains struct {
	DrainURLs []scalingo.LogDrain
	AppName   string
}

type ListAddonOpts struct {
	WithAddons bool
	AddonID    string
}

func List(ctx context.Context, app string, opts ListAddonOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	appToPrint := []printableDrains{}

	if opts.AddonID == "" {
		logDrains, err := c.LogDrainsList(ctx, app)
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
		addons, err := c.AddonsList(ctx, app)
		if err != nil {
			return errgo.Notef(err, "fail to list addons")
		}

		for _, addon := range addons {
			if opts.AddonID == addon.ID || opts.WithAddons {
				drains, err := c.LogDrainsAddonList(ctx, app, addon.ID)
				if err != nil {
					io.Status(err)
				}
				if len(drains) > 0 {
					appToPrint = append(appToPrint, printableDrains{
						AppName:   addon.AddonProvider.Name,
						DrainURLs: drains,
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
