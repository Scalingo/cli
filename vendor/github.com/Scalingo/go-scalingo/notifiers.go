package scalingo

import (
	"github.com/Scalingo/go-scalingo/debug"
	errgo "gopkg.in/errgo.v1"
)

type notifierCreateRequest struct {
	Notifier interface{} `json:"notifier"`
}

type NotifierCreateParams struct {
	Active         bool
	Name           string
	SendAllEvents  bool
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

func (c *Client) NotifierProvision(app, notifierType string, params NotifierCreateParams) (NotifierRes, error) {
	var notifierRes NotifierRes
	notifier := NewNotifier(notifierType, params)
	notifierParams := &notifierCreateRequest{Notifier: notifier}
	debug.Printf("[Notifier params]\n%+v", notifier)

	err := c.subresourceAdd(app, "notifiers", notifierParams, &notifierRes)
	if err != nil {
		return NotifierRes{}, errgo.Mask(err, errgo.Any)
	}
	return notifierRes, nil
}

func (c *Client) NotifierDestroy(app, ID string) error {
	return c.subresourceDelete(app, "notifiers", ID)
}
