package apps

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/events"
	"github.com/Scalingo/go-scalingo/v6"
)

func Events(ctx context.Context, app string, paginationOpts scalingo.PaginationOpts) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	appEvents, pagination, err := c.EventsList(ctx, app, paginationOpts)
	if err != nil {
		return errgo.Mask(err)
	}

	return events.DisplayTimeline(appEvents, pagination, events.DisplayTimelineOpts{})
}
