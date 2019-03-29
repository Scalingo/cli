package scalingo

import (
	"time"

	"gopkg.in/errgo.v1"
)

type AlertsService interface {
	AlertsList(app string) ([]*Alert, error)
	AlertAdd(app string, params AlertAddParams) (*Alert, error)
	AlertShow(app, id string) (*Alert, error)
	AlertUpdate(app, id string, params AlertUpdateParams) (*Alert, error)
	AlertRemove(app, id string) error
}

var _ AlertsService = (*Client)(nil)

// AlertInternal is an alert with the time.Duration attributes being the string
// representation. This is how the alerter service expects them.
type AlertInternal struct {
	ID            string  `json:"id"`
	AppID         string  `json:"app_id"`
	ContainerType string  `json:"container_type"`
	Metric        string  `json:"metric"`
	Limit         float64 `json:"limit"`
	Disabled      bool    `json:"disabled"`
	SendWhenBelow bool    `json:"send_when_below"`
	// DurationBeforeTrigger and RemindEvery are string parse-able using
	// time.ParseDuration
	DurationBeforeTrigger string `json:"duration_before_trigger"`
	RemindEvery           string `json:"remind_every"`
}

type AlertsRes struct {
	Alerts []*Alert `json:"alerts"`
}

type AlertRes struct {
	Alert *Alert `json:"alert"`
}

// Alert represents an alert with the duration attributes being...
// time.Duration.
type Alert struct {
	*AlertInternal
	DurationBeforeTrigger time.Duration `json:"duration_before_trigger"`
	RemindEvery           time.Duration `json:"remind_every"`
}

func (c *Client) AlertsList(app string) ([]*Alert, error) {
	var alertsRes AlertsRes
	err := c.ScalingoAPI().SubresourceList("apps", app, "alerts", nil, &alertsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to query the API to list an alert")
	}
	return alertsRes.Alerts, nil
}

type AlertAddParams struct {
	ContainerType         string
	Metric                string
	Limit                 float64
	RemindEvery           *time.Duration
	DurationBeforeTrigger *time.Duration
	SendWhenBelow         bool
	Notifiers             []string
}

func (c *Client) AlertAdd(app string, params AlertAddParams) (*Alert, error) {
	var alertRes AlertRes
	a := &AlertInternal{
		ContainerType: params.ContainerType,
		Metric:        params.Metric,
		Limit:         params.Limit,
		SendWhenBelow: params.SendWhenBelow,
	}
	if params.RemindEvery != nil {
		a.RemindEvery = (*params.RemindEvery).String()
	}
	if params.DurationBeforeTrigger != nil {
		a.DurationBeforeTrigger = (*params.DurationBeforeTrigger).String()
	}
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "alerts", struct {
		Alert *AlertInternal `json:"alert"`
	}{
		Alert: a,
	}, &alertRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to query the API to create an alert")
	}
	return alertRes.Alert, nil
}

func (c *Client) AlertShow(app, id string) (*Alert, error) {
	var alertRes AlertRes
	err := c.ScalingoAPI().SubresourceGet("apps", app, "alerts", id, nil, &alertRes)
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

func (c *Client) AlertUpdate(app, id string, params AlertUpdateParams) (*Alert, error) {
	var alertRes AlertRes
	err := c.ScalingoAPI().SubresourceUpdate("apps", app, "alerts", id, params, &alertRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to query the API to update an alert")
	}
	return alertRes.Alert, nil
}

func (c *Client) AlertRemove(app, id string) error {
	err := c.ScalingoAPI().SubresourceDelete("apps", app, "alerts", id)
	if err != nil {
		return errgo.Notef(err, "fail to query the API to remove an alert")
	}
	return nil
}
