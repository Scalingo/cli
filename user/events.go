package user

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/events"
	"github.com/Scalingo/go-utils/pagination"
)

func Events(ctx context.Context, paginationReq pagination.Request) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	appEvents, pagination, err := c.UserEventsList(ctx, paginationReq)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}

	return events.DisplayTimeline(appEvents, pagination, events.DisplayTimelineOpts{DisplayAppName: true})
}
