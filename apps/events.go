package apps

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/events"
	"github.com/Scalingo/go-utils/pagination"
)

func Events(ctx context.Context, app string, paginationReq pagination.Request) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	appEvents, pagination, err := c.EventsList(ctx, app, paginationReq)
	if err != nil {
		return errors.Wrapf(ctx, err, "list events for app %s", app)
	}

	return events.DisplayTimeline(appEvents, pagination, events.DisplayTimelineOpts{})
}
