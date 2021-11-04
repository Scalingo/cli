package log_drains

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

type RemoveAddonOpts struct {
	AddonID string
	OnlyApp bool
	URL     string
}

func Remove(app string, opts RemoveAddonOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to remove a log drain from the application")
	}

	if opts.AddonID != "" {
		// addon only
		err := c.LogDrainAddonRemove(app, opts.AddonID, opts.URL)
		if err != nil {
			return errgo.Notef(err, "fail to remove the log drain from the addon %s", opts.AddonID)
		}
		io.Status("The log drain", opts.URL, "has been deleted from the addon", opts.AddonID)
		return nil
	}

	err = c.LogDrainRemove(app, opts.URL)
	if err != nil {
		io.Status("fail to remove the log drain from the application:", app, "\n\t", err)
	} else {
		io.Status("Log drain", opts.URL, "has been deleted from the application", app)
	}

	if !opts.OnlyApp {
		addons, err := c.AddonsList(app)
		if err != nil {
			return errgo.Notef(err, "fail to list addons to remove log drain")
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

	return nil
}
