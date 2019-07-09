package scalingo

import (
	"encoding/json"
	"time"

	"github.com/Scalingo/go-scalingo/debug"
	errgo "gopkg.in/errgo.v1"
)

type NotifiersService interface {
	NotifiersList(app string) (Notifiers, error)
	NotifierProvision(app string, params NotifierParams) (*Notifier, error)
	NotifierByID(app, ID string) (*Notifier, error)
	NotifierUpdate(app, ID string, params NotifierParams) (*Notifier, error)
	NotifierDestroy(app, ID string) error
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

func (c *Client) NotifiersList(app string) (Notifiers, error) {
	var notifiersRes notifiersRequestRes
	err := c.ScalingoAPI().SubresourceList("apps", app, "notifiers", nil, &notifiersRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	var notifiers Notifiers
	for _, not := range notifiersRes.Notifiers {
		notifiers = append(notifiers, not.Specialize())
	}
	return notifiers, nil
}

func (c *Client) NotifierProvision(app string, params NotifierParams) (*Notifier, error) {
	var notifierRes notifierRequestRes
	notifier := newOutputNotifier(params)
	notifierRequestParams := &notifierRequestParams{notifier}

	err := c.ScalingoAPI().SubresourceAdd("apps", app, "notifiers", notifierRequestParams, &notifierRes)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierByID(app, ID string) (*Notifier, error) {
	var notifierRes notifierRequestRes
	err := c.ScalingoAPI().SubresourceGet("apps", app, "notifiers", ID, nil, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierUpdate(app, ID string, params NotifierParams) (*Notifier, error) {
	var notifierRes notifierRequestRes
	notifier := newOutputNotifier(params)
	notifierRequestParams := &notifierRequestParams{notifier}

	debug.Printf("[Notifier params]\n%+v", notifier)

	err := c.ScalingoAPI().SubresourceUpdate("apps", app, "notifiers", ID, notifierRequestParams, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierDestroy(app, ID string) error {
	return c.ScalingoAPI().SubresourceDelete("apps", app, "notifiers", ID)
}
