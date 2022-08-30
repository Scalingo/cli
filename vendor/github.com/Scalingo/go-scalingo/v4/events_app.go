package scalingo

import (
	"fmt"
	"strings"
)

type EventNewAppTypeData struct {
	GitSource string `json:"git_source"`
}

type EventNewAppType struct {
	Event
	TypeData EventNewAppTypeData `json:"type_data"`
}

func (ev *EventNewAppType) String() string {
	return "the application has been created"
}

type EventEditAppTypeData struct {
	ForceHTTPS *bool `json:"force_https"`
}

type EventEditAppType struct {
	Event
	TypeData EventEditAppTypeData `json:"type_data"`
}

func (ev *EventEditAppType) String() string {
	base := "application settings have been updated"
	if ev.TypeData.ForceHTTPS != nil {
		if *ev.TypeData.ForceHTTPS {
			base += ", Force HTTPS has been enabled"
		} else {
			base += ", Force HTTPS has been disabled"
		}
	}
	return base
}

type EventDeleteAppType struct {
	Event
}

func (ev *EventDeleteAppType) String() string {
	return "the application has been deleted"
}

type EventRenameAppTypeData struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type EventRenameAppType struct {
	Event
	TypeData EventRenameAppTypeData `json:"type_data"`
}

func (ev *EventRenameAppType) String() string {
	return fmt.Sprintf(
		"the application has been renamed from '%s' to '%s'",
		ev.TypeData.OldName, ev.TypeData.NewName,
	)
}

type EventTransferAppTypeData struct {
	OldOwner EventUser `json:"old_owner"`
	NewOwner EventUser `json:"new_owner"`
}

type EventTransferAppType struct {
	Event
	TypeData EventTransferAppTypeData `json:"type_data"`
}

func (ev *EventTransferAppType) String() string {
	return fmt.Sprintf(
		"the application has been transferred to %s (%s)",
		ev.TypeData.NewOwner.Username, ev.TypeData.NewOwner.Email,
	)
}

type EventRestartTypeData struct {
	Scope     []string `json:"scope"`
	AddonName string   `json:"addon_name"`
}

type EventRestartType struct {
	Event
	TypeData EventRestartTypeData `json:"type_data"`
}

func (ev *EventRestartType) String() string {
	if len(ev.TypeData.Scope) != 0 {
		return fmt.Sprintf("containers %v have been restarted", ev.TypeData.Scope)
	}
	return "containers have been restarted"
}

func (ev *EventRestartType) Who() string {
	if ev.TypeData.AddonName != "" {
		return fmt.Sprintf("Addon %s", ev.TypeData.AddonName)
	}
	return ev.Event.Who()
}

type EventStopAppTypeData struct {
	Reason string `json:"reason"`
}

type EventStopAppType struct {
	Event
	TypeData EventStopAppTypeData `json:"type_data"`
}

func (ev *EventStopAppType) String() string {
	return fmt.Sprintf("app has been stopped (reason: %s)", ev.TypeData.Reason)
}

func (ev *EventScaleType) String() string {
	return fmt.Sprintf(
		"containers have been scaled from %s, to %s",
		ev.TypeData.containersString(ev.TypeData.PreviousContainers),
		ev.TypeData.containersString(ev.TypeData.Containers),
	)
}

type EventScaleTypeData struct {
	PreviousContainers map[string]string `json:"previous_containers"`
	Containers         map[string]string `json:"containers"`
}

type EventScaleType struct {
	Event
	TypeData EventScaleTypeData `json:"type_data"`
}

func (e *EventScaleTypeData) containersString(containers map[string]string) string {
	types := []string{}
	for name, amountAndSize := range containers {
		types = append(types, fmt.Sprintf("%s:%s", name, amountAndSize))
	}
	return "[" + strings.Join(types, ", ") + "]"
}

type EventCrashTypeData struct {
	ContainerType string `json:"container_type"`
	CrashLogs     string `json:"crash_logs"`
	LogsURL       string `json:"logs_url"`
}

type EventCrashType struct {
	Event
	TypeData EventCrashTypeData `json:"type_data"`
}

func (ev *EventCrashType) String() string {
	msg := fmt.Sprintf("container '%v' has crashed", ev.TypeData.ContainerType)

	if ev.TypeData.CrashLogs != "" {
		msg += fmt.Sprintf(" (logs on %s)", ev.TypeData.LogsURL)
	}

	return msg
}

type EventRepeatedCrashTypeData struct {
	ContainerType string `json:"container_type"`
	CrashLogs     string `json:"crash_logs"`
	LogsURL       string `json:"logs_url"`
}

type EventRepeatedCrashType struct {
	Event
	TypeData EventRepeatedCrashTypeData `json:"type_data"`
}

func (ev *EventRepeatedCrashType) String() string {
	msg := fmt.Sprintf("container '%v' has crashed repeatedly", ev.TypeData.ContainerType)

	if ev.TypeData.CrashLogs != "" {
		msg += fmt.Sprintf(" (logs on %s)", ev.TypeData.LogsURL)
	}

	return msg
}

type EventDeploymentTypeData struct {
	DeploymentID   string `json:"deployment_id"`
	Pusher         string `json:"pusher"`
	GitRef         string `json:"git_ref"`
	Status         string `json:"status"`
	Duration       int    `json:"duration"`
	DeploymentUUID string `json:"deployment_uuid"`
}

type EventDeploymentType struct {
	Event
	TypeData EventDeploymentTypeData `json:"type_data"`
}

func (ev *EventDeploymentType) String() string {
	return fmt.Sprintf("deployment of %s (%s)", ev.TypeData.GitRef, ev.TypeData.Status)
}
