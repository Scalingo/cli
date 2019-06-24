package scalingo

import (
	"encoding/json"

	"github.com/Scalingo/go-scalingo/debug"
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
func (not *Notifier) GetNotifier() *Notifier {
	return not
}

func (not *Notifier) GetID() string {
	return not.ID
}

func (not *Notifier) GetName() string {
	return not.Name
}

func (not *Notifier) GetType() NotifierType {
	return not.Type
}

func (not *Notifier) GetSendAllEvents() bool {
	return *not.SendAllEvents
}

func (not *Notifier) GetSendAllAlerts() bool {
	return *not.SendAllAlerts
}

func (not *Notifier) GetSelectedEventIDs() []string {
	return not.SelectedEventIDs
}

func (not *Notifier) IsActive() bool {
	return *not.Active
}

func (not *Notifier) When() string {
	return not.UpdatedAt.Format("Mon Jan 02 2006 15:04:05")
}

func (not *Notifier) TypeDataPtr() interface{} {
	return &not.TypeData
}

func (not *Notifier) TypeDataMap() map[string]interface{} {
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

func (e *NotifierWebhookType) TypeDataPtr() interface{} {
	return &e.TypeData
}

func (not *NotifierWebhookType) TypeDataMap() map[string]interface{} {
	return map[string]interface{}{
		"webhook url": not.TypeData.WebhookURL,
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

func (e *NotifierSlackType) TypeDataPtr() interface{} {
	return &e.TypeData
}

func (not *NotifierSlackType) TypeDataMap() map[string]interface{} {
	return map[string]interface{}{
		"webhook url": not.TypeData.WebhookURL,
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

func (e *NotifierEmailType) TypeDataPtr() interface{} {
	return &e.TypeData
}

func (not *NotifierEmailType) TypeDataMap() map[string]interface{} {
	return map[string]interface{}{
		"emails":   not.TypeData.Emails,
		"user_ids": not.TypeData.UserIDs,
	}
}

func (pnot *Notifier) Specialize() DetailedNotifier {
	var detailedNotifier DetailedNotifier
	notifier := *pnot
	switch notifier.Type {
	case NotifierWebhook:
		detailedNotifier = &NotifierWebhookType{Notifier: notifier}
	case NotifierSlack:
		detailedNotifier = &NotifierSlackType{Notifier: notifier}
	case NotifierEmail:
		detailedNotifier = &NotifierEmailType{Notifier: notifier}
	default:
		return pnot
	}
	err := json.Unmarshal(pnot.RawTypeData, detailedNotifier.TypeDataPtr())
	if err != nil {
		debug.Printf("error reading the data: %+v\n", err)
		return pnot
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
