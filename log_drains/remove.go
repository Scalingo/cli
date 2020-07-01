package log_drains

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

type RemoveAddonOpts struct {
	WithAddons bool
	AddonID    string
	OnlyApp    bool
	URL        string
}

func Remove(app string, opts RemoveAddonOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to remove a log drain from the application")
	}

	if opts.OnlyApp {
		// app only
		err = c.LogDrainRemove(app, opts.URL)
		if err != nil {
			return errgo.Notef(err, "fail to remove the log drain from the application")
		}
	} else if opts.AddonID != "" {
		// addons only
		err := c.LogDrainAddonRemove(app, opts.AddonID, opts.URL)
		if err != nil {
			return errgo.Notef(err, "fail to remove the log drain from the addon %s", opts.AddonID)
		}
	} else {
		// app + addons
		addons, err := c.AddonsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list addons")
		}

		for _, addon := range addons {
			err := c.LogDrainAddonRemove(app, addon.ID, opts.URL)
			if err != nil {
				io.Status("fail to remove the log drain from the addon:", addon.AddonProvider.Name, "\n\t", err)
			} else {
				io.Status("Log drain", opts.URL, "has been deleted from the addon", addon.AddonProvider.Name)
			}
		}
	}

	io.Status("The log drain:", opts.URL, "has been deleted")
	return nil
}
