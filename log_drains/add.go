package log_drains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Add(app string, opts ListAddonOpts, params scalingo.LogDrainAddParams) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to add a log drain")
	}

	if opts.AddonID == "" || opts.WithAddons {
		d, err := c.LogDrainAdd(app, params)
		if err != nil {
			io.Status("fail to add drain to", "'"+app+"'", "application:\n\t", err)
			// return errgo.Notef(err, "fail to add drain to the application")
		} else {
			io.Status("Log drain", d.Drain.URL, "has been added to the application", app)
		}
	}

	addons, err := c.AddonsList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list addons")
	}

	for _, addon := range addons {
		if opts.AddonID == addon.ID || opts.WithAddons {
			// TODO(pc): do we need to test if the addon is a DB ?
			d, err := c.LogDrainAddonAdd(app, addon.ID, params)
			if err != nil {
				io.Status("fail to add drain to", "'"+addon.AddonProvider.Name+"'", "addon:\n\t", err)
				// return errgo.Notef(err, "fail to add drain to an addon")
			} else {
				io.Status("Log drain", d.Drain.URL, "has been added to the addon", addon.AddonProvider.Name)
			}

			if !opts.WithAddons {
				return nil
			}
		}
	}
	return nil
}
