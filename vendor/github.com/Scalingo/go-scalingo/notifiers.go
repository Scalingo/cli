package scalingo

import (
	"github.com/Scalingo/go-scalingo/debug"
	errgo "gopkg.in/errgo.v1"
)

type NotifiersService interface {
	NotifiersList(app string) (Notifiers, error)
	NotifierProvision(app, notifierType string, params NotifierParams) (*Notifier, error)
	NotifierByID(app, ID string) (*Notifier, error)
	NotifierUpdate(app, ID, notifierType string, params NotifierParams) (*Notifier, error)
	NotifierDestroy(app, ID string) error
}

var _ NotifiersService = (*Client)(nil)

// NotifierParams will be given as a parameter in notifiers function's
type NotifierParams struct {
	Active         *bool
	Name           string
	SendAllEvents  *bool
	SelectedEvents []string
	PlatformID     string

	// Options
	PhoneNumber string // SMS notifier
	Email       string // Email notifier
	WebhookURL  string // Webhook and Slack notifier
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
	err := c.subresourceList(app, "notifiers", nil, &notifiersRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	var notifiers Notifiers
	for _, not := range notifiersRes.Notifiers {
		notifiers = append(notifiers, not.Specialize())
	}
	return notifiers, nil
}

func (c *Client) NotifierProvision(app, notifierType string, params NotifierParams) (*Notifier, error) {
	var notifierRes notifierRequestRes
	notifier := NewOutputNotifier(notifierType, params)
	notifierRequestParams := &notifierRequestParams{notifier}

	err := c.subresourceAdd(app, "notifiers", notifierRequestParams, &notifierRes)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierByID(app, ID string) (*Notifier, error) {
	var notifierRes notifierRequestRes
	err := c.subresourceGet(app, "notifiers", ID, nil, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierUpdate(app, ID, notifierType string, params NotifierParams) (*Notifier, error) {
	var notifierRes notifierRequestRes
	notifier := NewOutputNotifier(notifierType, params)
	notifierRequestParams := &notifierRequestParams{notifier}

	debug.Printf("[Notifier params]\n%+v", notifier)

	err := c.subresourceUpdate(app, "notifiers", ID, notifierRequestParams, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierDestroy(app, ID string) error {
	return c.subresourceDelete(app, "notifiers", ID)
}
