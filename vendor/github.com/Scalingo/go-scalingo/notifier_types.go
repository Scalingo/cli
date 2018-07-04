package scalingo

import (
	"encoding/json"
	"time"

	"github.com/Scalingo/go-scalingo/debug"
)

// Used to omit attributes
type omit *struct{}

// Sruct used to represent a notifier.
type Notifier struct {
	ID             string                 `json:"id"`
	Active         *bool                  `json:"active,omitempty"`
	Name           string                 `json:"name,omitempty"`
	Type           NotifierType           `json:"type"`
	SendAllEvents  *bool                  `json:"send_all_events,omitempty"`
	SelectedEvents []EventTypeStruct      `json:"selected_events,omitempty"`
	TypeData       map[string]interface{} `json:"-"`
	RawTypeData    json.RawMessage        `json:"type_data"`
	PlatformID     string                 `json:"platform_id"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// Struct used to serialize a notifier
type NotifierOutput struct {
	*Notifier
	SelectedEvents []string    `json:"selected_events,omitempty"`
	TypeData       interface{} `json:"type_data,omitempty"`
	RawTypeData    omit        `json:",omitempty"` // Will always be empty and not serialized
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
	GetSelectedEvents() []EventTypeStruct
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

func (not *Notifier) GetSelectedEvents() []EventTypeStruct {
	return not.SelectedEvents
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

func NewOutputNotifier(notifierType string, params NotifierParams) NotifierOutput {
	detailedNotifier := NewDetailedNotifier(notifierType, params)
	res := NotifierOutput{
		Notifier:       detailedNotifier.GetNotifier(),
		TypeData:       detailedNotifier.TypeDataPtr(),
		SelectedEvents: params.SelectedEvents,
	}
	return res
}
