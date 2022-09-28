package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"
)

type AlertsService interface {
	AlertsList(ctx context.Context, app string) ([]*Alert, error)
	AlertAdd(ctx context.Context, app string, params AlertAddParams) (*Alert, error)
	AlertShow(ctx context.Context, app, id string) (*Alert, error)
	AlertUpdate(ctx context.Context, app, id string, params AlertUpdateParams) (*Alert, error)
	AlertRemove(ctx context.Context, app, id string) error
}

var _ AlertsService = (*Client)(nil)

type Alert struct {
	ID                    string                 `json:"id"`
	AppID                 string                 `json:"app_id"`
	ContainerType         string                 `json:"container_type"`
	Metric                string                 `json:"metric"`
	Limit                 float64                `json:"limit"`
	Disabled              bool                   `json:"disabled"`
	SendWhenBelow         bool                   `json:"send_when_below"`
	DurationBeforeTrigger time.Duration          `json:"duration_before_trigger"`
	RemindEvery           string                 `json:"remind_every"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	Metadata              map[string]interface{} `json:"metadata"`
	Notifiers             []string               `json:"notifiers"`
}

type AlertsRes struct {
	Alerts []*Alert `json:"alerts"`
}

type AlertRes struct {
	Alert *Alert `json:"alert"`
}

func (c *Client) AlertsList(ctx context.Context, app string) ([]*Alert, error) {
	var alertsRes AlertsRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "alerts", nil, &alertsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to query the API to list an alert")
	}
	return alertsRes.Alerts, nil
}

type AlertAddParams struct {
	ContainerType         string
	Metric                string
	Limit                 float64
	Disabled              bool
	RemindEvery           *time.Duration
	DurationBeforeTrigger *time.Duration
	SendWhenBelow         bool
	Notifiers             []string
}

func (c *Client) AlertAdd(ctx context.Context, app string, params AlertAddParams) (*Alert, error) {
	var alertRes AlertRes
	a := &Alert{
		ContainerType: params.ContainerType,
		Metric:        params.Metric,
		Limit:         params.Limit,
		SendWhenBelow: params.SendWhenBelow,
		Disabled:      params.Disabled,
		Notifiers:     params.Notifiers,
	}
	if params.RemindEvery != nil {
		a.RemindEvery = (*params.RemindEvery).String()
	}
	if params.DurationBeforeTrigger != nil {
		a.DurationBeforeTrigger = *params.DurationBeforeTrigger
	}
	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "alerts", AlertRes{
		Alert: a,
	}, &alertRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to query the API to create an alert")
	}
	return alertRes.Alert, nil
}

func (c *Client) AlertShow(ctx context.Context, app, id string) (*Alert, error) {
	var alertRes AlertRes
	err := c.ScalingoAPI().SubresourceGet(ctx, "apps", app, "alerts", id, nil, &alertRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to query the API to show an alert")
	}
	return alertRes.Alert, nil
}

type AlertUpdateParams struct {
	ContainerType         *string        `json:"container_type,omitempty"`
	Metric                *string        `json:"metric,omitempty"`
	Limit                 *float64       `json:"limit,omitempty"`
	Disabled              *bool          `json:"disabled,omitempty"`
	DurationBeforeTrigger *time.Duration `json:"duration_before_trigger,omitempty"`
	RemindEvery           *time.Duration `json:"remind_every,omitempty"`
	SendWhenBelow         *bool          `json:"send_when_below,omitempty"`
	Notifiers             *[]string      `json:"notifiers,omitempty"`
}

func (c *Client) AlertUpdate(ctx context.Context, app, id string, params AlertUpdateParams) (*Alert, error) {
	var alertRes AlertRes
	err := c.ScalingoAPI().SubresourceUpdate(ctx, "apps", app, "alerts", id, params, &alertRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to query the API to update an alert")
	}
	return alertRes.Alert, nil
}

func (c *Client) AlertRemove(ctx context.Context, app, id string) error {
	err := c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "alerts", id)
	if err != nil {
		return errgo.Notef(err, "fail to query the API to remove an alert")
	}
	return nil
}
