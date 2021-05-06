package log_drains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

type AddDrainOpts struct {
	WithAddons bool
	AddonID    string
	Params     scalingo.LogDrainAddParams
}

func Add(app string, opts AddDrainOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to add a log drain")
	}

	if opts.AddonID == "" || opts.WithAddons {
		d, err := c.LogDrainAdd(app, opts.Params)
		if err != nil {
			io.Status("fail to add drain to", "'"+app+"'", "application:\n\t", err)
		} else {
			io.Status("Log drain", d.Drain.URL, "has been added to the application", app)
		}
	}
	if !opts.WithAddons && opts.AddonID == "" {
		return nil
	}

	isAddonIDPresent := false
	addons, err := c.AddonsList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list addons")
	}

	for _, addon := range addons {
		if opts.AddonID == addon.ID {
			isAddonIDPresent = true
		}

		if opts.AddonID == addon.ID || opts.WithAddons {
			d, err := c.LogDrainAddonAdd(app, addon.ID, opts.Params)
			if err != nil {
				io.Status("fail to add drain to", "'"+addon.AddonProvider.Name+"'", "addon:\n\t", err)
			} else {
				io.Status("Log drain", d.Drain.URL, "has been added to the addon", addon.AddonProvider.Name)
			}

			if !opts.WithAddons {
				return nil
			}
		}
	}
	if !isAddonIDPresent && opts.AddonID != "" {
		return errgo.Notef(nil, "fail to add addon: addon_uuid doesn't exist for this application")
	}
	return nil
}
