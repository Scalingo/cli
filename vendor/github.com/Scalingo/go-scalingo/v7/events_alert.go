package scalingo

import "fmt"

type EventAlertTypeData struct {
	SenderName    string  `json:"sender_name"`
	Metric        string  `json:"metric"`
	Limit         float64 `json:"limit"`
	LimitText     string  `json:"limit_text"`
	Value         float64 `json:"value"`
	ValueText     string  `json:"value_text"`
	Activated     bool    `json:"activated"`
	SendWhenBelow bool    `json:"send_when_below"`
	RawLimit      float64 `json:"raw_limit"`
	RawValue      float64 `json:"raw_value"`
}

type EventAlertType struct {
	Event
	TypeData EventAlertTypeData `json:"type_data"`
}

func (ev *EventAlertType) String() string {
	d := ev.TypeData
	var res string
	if d.SendWhenBelow {
		res = "Low"
	} else {
		res = "High"
	}
	res += fmt.Sprintf(" %s usage on container %s ", d.Metric, d.SenderName)
	if ev.TypeData.Activated {
		res += "triggered"
	} else {
		res += "resolved"
	}
	res += fmt.Sprintf(" (value: %s, limit: %s)", d.ValueText, d.LimitText)

	return res
}

type EventNewAlertTypeData struct {
	ContainerType string  `json:"container_type"`
	Metric        string  `json:"metric"`
	Limit         float64 `json:"limit"`
	LimitText     string  `json:"limit_text"`
	SendWhenBelow bool    `json:"send_when_below"`
}

type EventNewAlertType struct {
	Event
	TypeData EventNewAlertTypeData `json:"type_data"`
}

func (ev *EventNewAlertType) String() string {
	d := ev.TypeData
	return fmt.Sprintf("Alert created about %s on container %s (limit: %s)", d.Metric, d.ContainerType, d.LimitText)
}

type EventDeleteAlertTypeData struct {
	ContainerType string `json:"container_type"`
	Metric        string `json:"metric"`
}

type EventDeleteAlertType struct {
	Event
	TypeData EventDeleteAlertTypeData `json:"type_data"`
}

func (ev *EventDeleteAlertType) String() string {
	d := ev.TypeData
	return fmt.Sprintf("Alert deleted about %s on container %s", d.Metric, d.ContainerType)
}
