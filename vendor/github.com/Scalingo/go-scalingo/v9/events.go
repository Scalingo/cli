package scalingo

//go:generate go run cmd/gen_events_boilerplate/main.go
//go:generate go run cmd/gen_events_specialize/main.go

import (
	"context"
	"encoding/json"

	"github.com/Scalingo/go-scalingo/v9/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type EventsService interface {
	EventTypesList(context.Context) ([]EventType, error)
	EventCategoriesList(context.Context) ([]EventCategory, error)
	EventsList(ctx context.Context, app string, opts PaginationOpts) (Events, PaginationMeta, error)
	UserEventsList(context.Context, PaginationOpts) (Events, PaginationMeta, error)
}

var _ EventsService = (*Client)(nil)

type EventsRes struct {
	Events []*Event `json:"events"`
	Meta   struct {
		PaginationMeta PaginationMeta `json:"pagination"`
	}
}

func (c *Client) EventsList(ctx context.Context, app string, opts PaginationOpts) (Events, PaginationMeta, error) {
	var eventsRes EventsRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "events", opts.ToMap(), &eventsRes)
	if err != nil {
		return nil, PaginationMeta{}, errors.Wrap(ctx, err, "list app events")
	}
	var events Events
	for _, ev := range eventsRes.Events {
		events = append(events, ev.Specialize())
	}
	return events, eventsRes.Meta.PaginationMeta, nil
}

func (c *Client) UserEventsList(ctx context.Context, opts PaginationOpts) (Events, PaginationMeta, error) {
	req := &http.APIRequest{
		Endpoint: "/events",
		Params:   opts.ToMap(),
	}

	var eventsRes EventsRes
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, PaginationMeta{}, errors.Wrap(ctx, err, "list user events")
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&eventsRes)
	if err != nil {
		return nil, PaginationMeta{}, errors.Wrap(ctx, err, "decode user events response")
	}

	var events Events
	for _, ev := range eventsRes.Events {
		events = append(events, ev.Specialize())
	}
	return events, eventsRes.Meta.PaginationMeta, nil
}
