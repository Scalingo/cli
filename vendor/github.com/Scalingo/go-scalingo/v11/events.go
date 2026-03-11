package scalingo

//go:generate go run cmd/gen_events_boilerplate/main.go
//go:generate go run cmd/gen_events_specialize/main.go

import (
	"context"

	"github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/pagination"
)

type EventsService interface {
	EventTypesList(context.Context) ([]EventType, error)
	EventCategoriesList(context.Context) ([]EventCategory, error)
	EventsList(ctx context.Context, app string, paginationReq pagination.Request) (Events, pagination.Meta, error)
	UserEventsList(ctx context.Context, paginationReq pagination.Request) (Events, pagination.Meta, error)
}

var _ EventsService = (*Client)(nil)

type EventsRes struct {
	Events []*Event `json:"events"`
	Meta   struct {
		Pagination pagination.Meta `json:"pagination"`
	}
}

func (c *Client) EventsList(ctx context.Context, app string, paginationReq pagination.Request) (Events, pagination.Meta, error) {
	var eventsRes EventsRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "events", paginationReq.ToURLValues(), &eventsRes)
	if err != nil {
		return nil, pagination.Meta{}, errors.Wrap(ctx, err, "list app events")
	}
	var events Events
	for _, ev := range eventsRes.Events {
		events = append(events, ev.Specialize())
	}
	return events, eventsRes.Meta.Pagination, nil
}

func (c *Client) UserEventsList(ctx context.Context, paginationReq pagination.Request) (Events, pagination.Meta, error) {
	req := &http.APIRequest{
		Endpoint: "/events",
		Params:   paginationReq.ToURLValues(),
	}

	var eventsRes EventsRes
	err := c.ScalingoAPI().DoRequest(ctx, req, &eventsRes)
	if err != nil {
		return nil, pagination.Meta{}, errors.Wrap(ctx, err, "list user events")
	}

	var events Events
	for _, ev := range eventsRes.Events {
		events = append(events, ev.Specialize())
	}
	return events, eventsRes.Meta.Pagination, nil
}
