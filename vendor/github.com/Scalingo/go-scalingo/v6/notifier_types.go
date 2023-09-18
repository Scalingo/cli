package scalingo

import (
	"encoding/json"

	"github.com/Scalingo/go-scalingo/v6/debug"
)

// Used to omit attributes
type omit *struct{}

// Struct used to serialize a notifier
type NotifierOutput struct {
	*Notifier
	TypeData    NotifierTypeDataParams `json:"type_data,omitempty"`
	RawTypeData omit                   `json:",omitempty"` // Will always be empty and not serialized
}

type NotifierTypeDataParams struct {
	WebhookURL  string   `json:"webhook_url,omitempty"`
	Emails      []string `json:"emails,omitempty"`
	UserIDs     []string `json:"user_ids,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
}

type NotifierType string

const (
	NotifierWebhook NotifierType = "webhook"
	NotifierSlack   NotifierType = "slack"
	NotifierEmail   NotifierType = "email"
)

type DetailedNotifier interface {
	GetNotifier() *Notifier
	GetID() string
	GetName() string
	GetType() NotifierType
	GetSendAllEvents() bool
	GetSendAllAlerts() bool
	GetSelectedEventIDs() []string
	IsActive() bool
	When() string
	TypeDataPtr() interface{}
	TypeDataMap() map[string]interface{}
}

type Notifiers []DetailedNotifier

// DetailedNotifier implementation
func (n *Notifier) GetNotifier() *Notifier {
	return n
}

func (n *Notifier) GetID() string {
	return n.ID
}

func (n *Notifier) GetName() string {
	return n.Name
}

func (n *Notifier) GetType() NotifierType {
	return n.Type
}

func (n *Notifier) GetSendAllEvents() bool {
	return *n.SendAllEvents
}

func (n *Notifier) GetSendAllAlerts() bool {
	return *n.SendAllAlerts
}

func (n *Notifier) GetSelectedEventIDs() []string {
	return n.SelectedEventIDs
}

func (n *Notifier) IsActive() bool {
	return *n.Active
}

func (n *Notifier) When() string {
	return n.UpdatedAt.Format("Mon Jan 02 2006 15:04:05")
}

func (n *Notifier) TypeDataPtr() interface{} {
	return &n.TypeData
}

func (n *Notifier) TypeDataMap() map[string]interface{} {
	return map[string]interface{}{}
}

// Webhook
type NotifierWebhookType struct {
	Notifier
	TypeData NotifierWebhookTypeData `json:"type_data,omitempty"`
}

type NotifierWebhookTypeData struct {
	WebhookURL string `json:"webhook_url,omitempty"`
}

func (n *NotifierWebhookType) TypeDataPtr() interface{} {
	return &n.TypeData
}

func (n *NotifierWebhookType) TypeDataMap() map[string]interface{} {
	return map[string]interface{}{
		"webhook url": n.TypeData.WebhookURL,
	}
}

// Slack
type NotifierSlackType struct {
	Notifier
	TypeData NotifierSlackTypeData `json:"type_data,omitempty"`
}

type NotifierSlackTypeData struct {
	WebhookURL string `json:"webhook_url,omitempty"`
}

func (n *NotifierSlackType) TypeDataPtr() interface{} {
	return &n.TypeData
}

func (n *NotifierSlackType) TypeDataMap() map[string]interface{} {
	return map[string]interface{}{
		"webhook url": n.TypeData.WebhookURL,
	}
}

// Email
type NotifierEmailType struct {
	Notifier
	TypeData NotifierEmailTypeData `json:"type_data,omitempty"`
}

type NotifierEmailTypeData struct {
	Emails  []string `json:"emails,omitempty"`
	UserIDs []string `json:"user_ids,omitempty"`
}

func (n *NotifierEmailType) TypeDataPtr() interface{} {
	return &n.TypeData
}

func (n *NotifierEmailType) TypeDataMap() map[string]interface{} {
	return map[string]interface{}{
		"emails":   n.TypeData.Emails,
		"user_ids": n.TypeData.UserIDs,
	}
}

func (n *Notifier) Specialize() DetailedNotifier {
	var detailedNotifier DetailedNotifier
	notifier := *n
	switch notifier.Type {
	case NotifierWebhook:
		detailedNotifier = &NotifierWebhookType{Notifier: notifier}
	case NotifierSlack:
		detailedNotifier = &NotifierSlackType{Notifier: notifier}
	case NotifierEmail:
		detailedNotifier = &NotifierEmailType{Notifier: notifier}
	default:
		return n
	}
	err := json.Unmarshal(n.RawTypeData, detailedNotifier.TypeDataPtr())
	if err != nil {
		debug.Printf("error reading the data: %+v\n", err)
		return n
	}
	return detailedNotifier
}

func NewDetailedNotifier(notifierType string, params NotifierParams) DetailedNotifier {
	debug.Printf("[NewNotifier] notifierType: %+v\nparams: %+v\n", notifierType, params)
	var specializedNotifier DetailedNotifier
	notifier := &Notifier{
		Active:        params.Active,
		Name:          params.Name,
		SendAllEvents: params.SendAllEvents,
		PlatformID:    params.PlatformID,
	}

	switch notifierType {
	case "webhook":
		specializedNotifier = &NotifierWebhookType{
			Notifier: *notifier,
			TypeData: NotifierWebhookTypeData{
				WebhookURL: params.WebhookURL,
			},
		}
	case "slack":
		specializedNotifier = &NotifierSlackType{
			Notifier: *notifier,
			TypeData: NotifierSlackTypeData{
				WebhookURL: params.WebhookURL,
			},
		}
	case "email":
		specializedNotifier = &NotifierEmailType{
			Notifier: *notifier,
			TypeData: NotifierEmailTypeData{
				Emails:  params.Emails,
				UserIDs: params.UserIDs,
			},
		}
	}

	debug.Printf("[NewNotifier] result: %+v\n", specializedNotifier)
	return specializedNotifier
}

// newOutputNotifier prepares the payload to send to the API to
// create or update a notifier
func newOutputNotifier(params NotifierParams) NotifierOutput {
	res := NotifierOutput{
		Notifier: &Notifier{
			Active:           params.Active,
			Name:             params.Name,
			PlatformID:       params.PlatformID,
			SendAllAlerts:    params.SendAllAlerts,
			SendAllEvents:    params.SendAllEvents,
			SelectedEventIDs: params.SelectedEventIDs,
		},
		TypeData: NotifierTypeDataParams{
			Emails:      params.Emails,
			UserIDs:     params.UserIDs,
			WebhookURL:  params.WebhookURL,
			PhoneNumber: params.PhoneNumber,
		},
	}
	return res
}
