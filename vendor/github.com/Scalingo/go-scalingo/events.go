package scalingo

import (
	"encoding/json"

	"github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type EventsService interface {
	EventsList(app string, opts PaginationOpts) (Events, PaginationMeta, error)
	UserEventsList(opts PaginationOpts) (Events, PaginationMeta, error)
}

var _ EventsService = (*Client)(nil)

type EventsRes struct {
	Events []*Event `json:"events"`
	Meta   struct {
		PaginationMeta PaginationMeta `json:"pagination"`
	}
}

func (c *Client) EventsList(app string, opts PaginationOpts) (Events, PaginationMeta, error) {
	var eventsRes EventsRes
	err := c.ScalingoAPI().SubresourceList("apps", app, "events", opts.ToMap(), &eventsRes)
	if err != nil {
		return nil, PaginationMeta{}, errgo.Mask(err)
	}
	var events Events
	for _, ev := range eventsRes.Events {
		events = append(events, ev.Specialize())
	}
	return events, eventsRes.Meta.PaginationMeta, nil
}

func (c *Client) UserEventsList(opts PaginationOpts) (Events, PaginationMeta, error) {
	req := &http.APIRequest{
		Endpoint: "/events",
		Params:   opts.ToMap(),
	}

	var eventsRes EventsRes
	res, err := c.ScalingoAPI().Do(req)
	if err != nil {
		return nil, PaginationMeta{}, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&eventsRes)
	if err != nil {
		return nil, PaginationMeta{}, errgo.Mask(err, errgo.Any)
	}

	var events Events
	for _, ev := range eventsRes.Events {
		events = append(events, ev.Specialize())
	}
	return events, eventsRes.Meta.PaginationMeta, nil
}
