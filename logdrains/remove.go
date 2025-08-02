package logdrains

import (
	"context"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

type RemoveAddonOpts struct {
	AddonID string
	OnlyApp bool
	URL     string
}

func Remove(ctx context.Context, app string, opts RemoveAddonOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	if opts.AddonID != "" {
		// addon only
		err := c.LogDrainAddonRemove(ctx, app, opts.AddonID, opts.URL)
		if err != nil {
			return errors.Wrap(ctx, err, "remove log drain from addon "+opts.AddonID)
		}
		io.Status("The log drain", opts.URL, "has been deleted from the addon", opts.AddonID)
		return nil
	}

	err = c.LogDrainRemove(ctx, app, opts.URL)
	if err != nil {
		return errors.Wrap(ctx, err, "remove log drain from application "+app)
	}
	io.Status("Log drain", opts.URL, "has been deleted from the application", app)

	if !opts.OnlyApp {
		addons, err := c.AddonsList(ctx, app)
		if err != nil {
			return errors.Wrap(ctx, err, "list addons to remove log drain")
		}

		for _, addon := range addons {
			err := c.LogDrainAddonRemove(ctx, app, addon.ID, opts.URL)
			if err != nil {
				// Check if this is a "not found" error, which can happen if the log drain
				// was already removed by the main API call
				if strings.Contains(err.Error(), "not found") {
					io.Status("Log drain", opts.URL, "was already removed from the addon", addon.AddonProvider.Name)
				} else {
					io.Status("Unable to remove the log drain from the addon:", addon.AddonProvider.Name, "\n\t", err)
				}
			} else {
				io.Status("Log drain", opts.URL, "has been deleted from the addon", addon.AddonProvider.Name)
			}
		}
	}

	return nil
}
