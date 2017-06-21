package scalingo

import (
	"github.com/Scalingo/go-scalingo/debug"
	errgo "gopkg.in/errgo.v1"
)

type notifierRequest struct {
	Notifier interface{} `json:"notifier"`
}

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

type NotifierRes struct {
	Notifier Notifier `json:"notifier"`
}

type NotifiersRes struct {
	Notifiers []*Notifier `json:"notifiers"`
}

func (c *Client) NotifiersList(app string) (Notifiers, error) {
	var notifiersRes NotifiersRes
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
	var notifierRes NotifierRes
	notifier := NewNotifier(notifierType, params)
	notifierParams := &notifierRequest{Notifier: notifier}
	debug.Printf("[Notifier params]\n%+v", notifier)

	err := c.subresourceAdd(app, "notifiers", notifierParams, &notifierRes)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierByID(app, ID string) (*Notifier, error) {
	var notifierRes NotifierRes
	err := c.subresourceGet(app, "notifiers", ID, nil, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierUpdate(app, ID, notifierType string, params NotifierParams) (*Notifier, error) {
	var notifierRes NotifierRes
	notifier := NewNotifier(notifierType, params)
	notifierParams := &notifierRequest{Notifier: notifier}
	debug.Printf("[Notifier params]\n%+v", notifier)

	err := c.subresourceUpdate(app, "notifiers", ID, notifierParams, &notifierRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return &notifierRes.Notifier, nil
}

func (c *Client) NotifierDestroy(app, ID string) error {
	return c.subresourceDelete(app, "notifiers", ID)
}
