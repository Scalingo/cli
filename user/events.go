package user

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/events"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func Events(paginationOpts scalingo.PaginationOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	appEvents, pagination, err := c.UserEventsList(paginationOpts)
	if err != nil {
		return errgo.Mask(err)
	}

	return events.DisplayTimeline(appEvents, pagination, events.DisplayTimelineOpts{DisplayAppName: true})
}
