package apps

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/events"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func Events(app string, paginationOpts scalingo.PaginationOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	appEvents, pagination, err := c.EventsList(app, paginationOpts)
	if err != nil {
		return errgo.Mask(err)
	}

	return events.DisplayTimeline(appEvents, pagination, events.DisplayTimelineOpts{})
}
