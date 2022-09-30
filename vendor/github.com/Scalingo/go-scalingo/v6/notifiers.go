package scalingo

import (
	"context"
	"encoding/json"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v6/debug"
)

type NotifiersService interface {
	NotifiersList(ctx context.Context, app string) (Notifiers, error)
	NotifierProvision(ctx context.Context, app string, params NotifierParams) (*Notifier, error)
	NotifierByID(ctx context.Context, app, ID string) (*Notifier, error)
	NotifierUpdate(ctx context.Context, app, ID string, params NotifierParams) (*Notifier, error)
	NotifierDestroy(ctx context.Context, app, ID string) error
}

var _ NotifiersService = (*Client)(nil)

// Struct used to represent a notifier.
type Notifier struct {
	ID               string                 `json:"id"`
	AppID            string                 `json:"app_id"`
	Active           *bool                  `json:"active,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Type             NotifierType           `json:"type"`
	SendAllEvents    *bool                  `json:"send_all_events,omitempty"`
	SendAllAlerts    *bool                  `json:"send_all_alerts,omitempty"`
	SelectedEventIDs []string               `json:"selected_event_ids,omitempty"`
	TypeData         map[string]interface{} `json:"-"`
	RawTypeData      json.RawMessage        `json:"type_data"`
	PlatformID       string                 `json:"platform_id"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// NotifierParams will be given as a parameter in notifiers function's
type NotifierParams struct {
	Active           *bool
	Name             string
	SendAllEvents    *bool
	SendAllAlerts    *bool
	SelectedEventIDs []string
	PlatformID       string

	// Options
	PhoneNumber string   // SMS notifier
	Emails      []string // Email notifier
	UserIDs     []string // Email notifier
	WebhookURL  string   // Webhook and Slack notifier
}

// The struct that will be serialized in all notifier's request
type notifierRequestParams struct {
	NotifierOutput `json:"notifier"`
}

// The struct that will be deserialized from all notifier request response
type notifierRequestRes struct {
	Notifier Notifier `json:"notifier"`
}

type notifiersRequestRes struct {
	Notifiers []*Notifier `json:"notifiers"`
}

func (c *Client) NotifiersList(ctx context.Context, app string) (Notifiers, error) {
	var notifiersRes notifiersRequestRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "notifiers", nil, &notifiersRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	var notifiers Notifiers
	for _, not := range notifiersRes.Notifiers {
		notifiers = append(notifiers, not.Specialize())
	}
	return notifiers, nil
}

func (c *Client) NotifierProvision(ctx context.Context, app string, params NotifierParams) (*Notifier, error) {
	var notifierRes notifierRequestRes
	notifier := newOutputNotifier(params)
	notifierRequestParams := &notifierRequestParams{notifier}

	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "notifiers", notifierRequestParams, &notifierRes)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierByID(ctx context.Context, app, ID string) (*Notifier, error) {
	var notifierRes notifierRequestRes
	err := c.ScalingoAPI().SubresourceGet(ctx, "apps", app, "notifiers", ID, nil, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierUpdate(ctx context.Context, app, ID string, params NotifierParams) (*Notifier, error) {
	var notifierRes notifierRequestRes
	notifier := newOutputNotifier(params)
	notifierRequestParams := &notifierRequestParams{notifier}

	debug.Printf("[Notifier params]\n%+v", notifier)

	err := c.ScalingoAPI().SubresourceUpdate(ctx, "apps", app, "notifiers", ID, notifierRequestParams, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierDestroy(ctx context.Context, app, ID string) error {
	return c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "notifiers", ID)
}
