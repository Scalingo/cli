package db

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/logs"
)

type LogsOpts struct {
	Follow bool
	Count  int
}

// Logs displays the addon logs.
// app may be an app UUID or name.
// addon may be a addon UUID or an addon type (e.g. MongoDB).
func Logs(ctx context.Context, app, addonUUID string, opts LogsOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	url, err := c.AddonLogsURL(ctx, app, addonUUID)
	if err != nil {
		return errgo.Notef(err, "fail to get log URL")
	}

	err = logs.Dump(ctx, url, opts.Count, "")
	if err != nil {
		return errgo.Notef(err, "fail to dump logs")
	}

	if opts.Follow {
		err := logs.Stream(ctx, url, "")
		if err != nil {
			return errgo.Notef(err, "fail to stream logs")
		}
	}
	return nil
}
