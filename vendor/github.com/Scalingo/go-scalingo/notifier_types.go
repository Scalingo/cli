package scalingo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Scalingo/go-scalingo/debug"
)

type Notifier struct {
	ID             string                 `json:"id"`
	Active         bool                   `json:"active"`
	Name           string                 `json:"name"`
	Type           NotifierType           `json:"type"`
	SendAllEvents  bool                   `json:"send_all_events"`
	SelectedEvents []string               `json:"selected_events"`
	TypeData       map[string]interface{} `json:"-"`
	RawTypeData    json.RawMessage        `json:"type_data"`
	PlatformID     string                 `json:"platform_id"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type NotifierType string

const (
	NotifierWebhook NotifierType = "webhook"
	NotifierSlack                = "slack"
)

type DetailedNotifier interface {
	GetNotifier() *Notifier
	GetID() string
	GetName() string
	GetType() NotifierType
	GetSendAllEvents() bool
	GetSelectedEvents() []string
	IsActive() bool
	When() string
	TypeDataPtr() interface{}
	TypeDataString() string
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
	return not.SendAllEvents
}

func (not *Notifier) GetSelectedEvents() []string {
	return not.SelectedEvents
}

func (not *Notifier) IsActive() bool {
	return not.Active
}

func (not *Notifier) When() string {
	return not.UpdatedAt.Format("Mon Jan 02 2006 15:04:05")
}

func (not *Notifier) TypeDataPtr() interface{} {
	return &not.TypeData
}

func (not *Notifier) TypeDataString() string {
	return "unknow notifier type"
}

// Webhook
type NotifierWebhookType struct {
	Notifier
	TypeData NotifierWebhookTypeData `json:"type_data"`
}

type NotifierWebhookTypeData struct {
	WebhookURL string `json:"webhook_url"`
}

func (e *NotifierWebhookType) TypeDataPtr() interface{} {
	return &e.TypeData
}

func (not *NotifierWebhookType) TypeDataString() string {
	return fmt.Sprintf("- webhook url: %s", not.TypeData.WebhookURL)
}

// Slack
type NotifierSlackType struct {
	Notifier
	TypeData NotifierSlackTypeData `json:"type_data"`
}

type NotifierSlackTypeData struct {
	WebhookURL string `json:"webhook_url"`
}

func (e *NotifierSlackType) TypeDataPtr() interface{} {
	return &e.TypeData
}

func (not *NotifierSlackType) TypeDataString() string {
	return fmt.Sprintf("- webhook url: %s", not.TypeData.WebhookURL)
}

func (pnot *Notifier) Specialize() DetailedNotifier {
	var detailedNotifier DetailedNotifier
	notifier := *pnot
	switch notifier.Type {
	case NotifierWebhook:
		detailedNotifier = &NotifierWebhookType{Notifier: notifier}
	case NotifierSlack:
		detailedNotifier = &NotifierSlackType{Notifier: notifier}
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

func NewNotifier(notifierType string, params NotifierCreateParams) DetailedNotifier {
	debug.Printf("[NewNotifier] notifierType: %+v\nparams: %+v\n", notifierType, params)
	var specializedNotifier DetailedNotifier
	notifier := &Notifier{
		Active:         params.Active,
		Name:           params.Name,
		SendAllEvents:  params.SendAllEvents,
		SelectedEvents: params.SelectedEvents,
		PlatformID:     params.PlatformID,
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
	}

	debug.Printf("[NewNotifier] result: %+v\n", specializedNotifier)
	return specializedNotifier
}
