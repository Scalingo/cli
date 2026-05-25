package scalingo

import (
	"fmt"
	"strings"
	"time"
)

type EventAddon struct {
	AddonProviderName string `json:"addon_provider_name"`
	PlanName          string `json:"plan_name"`
	ResourceID        string `json:"resource_id"`
}

type EventNewAddonTypeData struct {
	EventAddon
}

type EventNewAddonType struct {
	Event

	TypeData EventNewAddonTypeData `json:"type_data"`
}

func (ev *EventNewAddonType) String() string {
	return fmt.Sprintf(
		"'%s' (%s) has been added (plan '%s')",
		ev.TypeData.ResourceID, ev.TypeData.AddonProviderName, ev.TypeData.PlanName,
	)
}

type EventUpgradeAddonTypeData struct {
	EventAddon

	OldPlanName string `json:"old_plan_name"`
	NewPlanName string `json:"new_plan_name"`
}

type EventUpgradeAddonType struct {
	Event

	TypeData EventUpgradeAddonTypeData `json:"type_data"`
}

func (ev *EventUpgradeAddonType) String() string {
	return fmt.Sprintf(
		"'%s' (%s) plan has been changed from '%s' to '%s'",
		ev.TypeData.ResourceID, ev.TypeData.AddonProviderName, ev.TypeData.OldPlanName, ev.TypeData.NewPlanName,
	)
}

type EventAddonUpdatedTypeData struct {
	AddonID           string `json:"addon_id"`
	AddonPlanName     string `json:"addon_plan_name"`
	AddonResourceID   string `json:"addon_resource_id"`
	AddonProviderID   string `json:"addon_provider_id"`
	AddonProviderName string `json:"addon_provider_name"`

	// Status has only two items when is updated, the old value and the new value, in this order
	Status []AddonStatus `json:"status"`
	// AttributesChanged contain names of changed attributes
	AttributesChanged []string `json:"attributes_changed"`
}

type EventAddonUpdatedType struct {
	Event

	TypeData EventAddonUpdatedTypeData `json:"type_data"`
}

func (ev *EventAddonUpdatedType) String() string {
	d := ev.TypeData
	return fmt.Sprintf(
		"Addon %s %s updated, status %v -> %v",
		d.AddonProviderName, d.AddonResourceID, d.Status[0], d.Status[1],
	)
}

type EventDeleteAddonTypeData struct {
	EventAddon
}

type EventDeleteAddonType struct {
	Event

	TypeData EventDeleteAddonTypeData `json:"type_data"`
}

func (ev *EventDeleteAddonType) String() string {
	return fmt.Sprintf(
		"'%s' (%s) plan has been deleted",
		ev.TypeData.ResourceID, ev.TypeData.AddonProviderName,
	)
}

type EventResumeAddonTypeData struct {
	EventAddon
}

type EventResumeAddonType struct {
	Event

	TypeData EventResumeAddonTypeData `json:"type_data"`
}

func (ev *EventResumeAddonType) String() string {
	return fmt.Sprintf(
		"'%s' (%s) has been resumed",
		ev.TypeData.ResourceID, ev.TypeData.AddonProviderName,
	)
}

type EventSuspendAddonTypeData struct {
	EventAddon

	Reason string `json:"reason"`
}

type EventSuspendAddonType struct {
	Event

	TypeData EventSuspendAddonTypeData `json:"type_data"`
}

func (ev *EventSuspendAddonType) String() string {
	return fmt.Sprintf(
		"'%s' (%s) has been suspended (reason: %s)",
		ev.TypeData.ResourceID, ev.TypeData.AddonProviderName, ev.TypeData.Reason,
	)
}

type EventDatabaseAddFeatureType struct {
	Event

	TypeData EventDatabaseAddFeatureTypeData `json:"type_data"`
}

type EventDatabaseAddFeatureTypeData struct {
	EventSecurityTypeData

	Feature           string `json:"feature"`
	AddonProviderID   string `json:"addon_provider_id"`
	AddonProviderName string `json:"addon_provider_name"`
	AddonUUID         string `json:"addon_uuid"`
}

func (ev *EventDatabaseAddFeatureType) String() string {
	return fmt.Sprintf(
		"Feature %s enabled for addon '%s' (%s) ",
		ev.TypeData.Feature, ev.TypeData.AddonUUID, ev.TypeData.AddonProviderName,
	)
}

type EventDatabaseRemoveFeatureType struct {
	Event

	TypeData EventDatabaseRemoveFeatureTypeData `json:"type_data"`
}

type EventDatabaseRemoveFeatureTypeData struct {
	EventSecurityTypeData

	Feature           string `json:"feature"`
	AddonProviderID   string `json:"addon_provider_id"`
	AddonProviderName string `json:"addon_provider_name"`
	AddonUUID         string `json:"addon_uuid"`
}

func (ev *EventDatabaseRemoveFeatureType) String() string {
	return fmt.Sprintf(
		"Feature %s disabled for addon '%s' (%s) ",
		ev.TypeData.Feature, ev.TypeData.AddonUUID, ev.TypeData.AddonProviderName,
	)
}

type EventDatabaseBackupSucceededType struct {
	Event

	TypeData EventDatabaseBackupSucceededTypeData `json:"type_data"`
}

func (ev *EventDatabaseBackupSucceededType) Who() string {
	if ev.User.Email == ScalingoDeployUserEmail {
		return "Scalingo Automated Backup Service"
	}

	return ev.Event.Who()
}

func (ev *EventDatabaseBackupSucceededType) String() string {
	methodStr := "B"

	switch ev.TypeData.BackupMethod {
	case BackupMethodPeriodic:
		methodStr = "Periodic b"
	case BackupMethodManual:
		methodStr = "Manual b"
	}
	return fmt.Sprintf(
		"%sackup %s for addon '%s' (%s) succeeded",
		methodStr,
		ev.TypeData.BackupID, ev.TypeData.AddonName, ev.TypeData.ResourceID,
	)
}

type EventDatabaseBackupSucceededTypeData struct {
	EventSecurityTypeData

	AddonUUID    string       `json:"addon_uuid"`
	AddonName    string       `json:"addon_name"`
	ResourceID   string       `json:"resource_id"`
	BackupMethod BackupMethod `json:"backup_method"`
	BackupID     string       `json:"backup_id"`
	BackupStatus string       `json:"backup_status"`
	StartedAt    time.Time    `json:"started_at"`
	EndedAt      time.Time    `json:"ended_at"`
}

type EventDatabaseBackupFailedType struct {
	Event

	TypeData EventDatabaseBackupFailedTypeData `json:"type_data"`
}

func (ev *EventDatabaseBackupFailedType) Who() string {
	if ev.User.Email == ScalingoDeployUserEmail {
		return "Scalingo Automated Backup Service"
	}

	return ev.Event.Who()
}

func (ev *EventDatabaseBackupFailedType) String() string {
	methodStr := ""

	switch ev.TypeData.BackupMethod {
	case BackupMethodPeriodic:
		methodStr = "Periodic"
	case BackupMethodManual:
		methodStr = "Manual"
	}

	return fmt.Sprintf(
		"%s backup %s for addon '%s' (%s) failed",
		methodStr,
		ev.TypeData.BackupID, ev.TypeData.AddonName, ev.TypeData.ResourceID,
	)
}

type EventDatabaseBackupFailedTypeData struct {
	AddonUUID    string       `json:"addon_uuid"`
	AddonName    string       `json:"addon_name"`
	BackupMethod BackupMethod `json:"backup_method"`
	ResourceID   string       `json:"resource_id"`
	BackupID     string       `json:"backup_id"`
	BackupStatus string       `json:"backup_status"`
	StartedAt    time.Time    `json:"started_at"`
	EndedAt      time.Time    `json:"ended_at"`
	EventSecurityTypeData
}

type EventDatabaseContinuousBackupTypeData struct {
	EventSecurityTypeData

	AddonName                    string    `json:"addon_name"`
	ResourceID                   string    `json:"resource_id"`
	AddonUUID                    string    `json:"addon_uuid"`
	Status                       string    `json:"status"`
	Error                        string    `json:"error"`
	Recoverable                  bool      `json:"recoverable"`
	CheckedAt                    time.Time `json:"checked_at"`
	UnrecoverableDurationSeconds int64     `json:"unrecoverable_duration_seconds"`
}

type EventDatabaseContinuousBackupHealthyType struct {
	Event

	TypeData EventDatabaseContinuousBackupHealthyTypeData `json:"type_data"`
}

func (ev *EventDatabaseContinuousBackupHealthyType) String() string {
	return formatDatabaseContinuousBackupString(
		ev.TypeData.EventDatabaseContinuousBackupTypeData,
		"healthy",
	)
}

func (ev *EventDatabaseContinuousBackupHealthyType) Who() string {
	return ev.Event.Who()
}

type EventDatabaseContinuousBackupHealthyTypeData struct {
	EventDatabaseContinuousBackupTypeData
}

type EventDatabaseContinuousBackupDelayedType struct {
	Event

	TypeData EventDatabaseContinuousBackupDelayedTypeData `json:"type_data"`
}

func (ev *EventDatabaseContinuousBackupDelayedType) String() string {
	return formatDatabaseContinuousBackupString(
		ev.TypeData.EventDatabaseContinuousBackupTypeData,
		"delayed",
	)
}

func (ev *EventDatabaseContinuousBackupDelayedType) Who() string {
	return ev.Event.Who()
}

type EventDatabaseContinuousBackupDelayedTypeData struct {
	EventDatabaseContinuousBackupTypeData
}

type EventDatabaseContinuousBackupStaleType struct {
	Event

	TypeData EventDatabaseContinuousBackupStaleTypeData `json:"type_data"`
}

func (ev *EventDatabaseContinuousBackupStaleType) String() string {
	return formatDatabaseContinuousBackupString(
		ev.TypeData.EventDatabaseContinuousBackupTypeData,
		"stale",
	)
}

func (ev *EventDatabaseContinuousBackupStaleType) Who() string {
	return ev.Event.Who()
}

type EventDatabaseContinuousBackupStaleTypeData struct {
	EventDatabaseContinuousBackupTypeData
}

func formatDatabaseContinuousBackupString(data EventDatabaseContinuousBackupTypeData, status string) string {
	message := fmt.Sprintf(
		"Point-in-time recovery for database '%s' is %s",
		data.ResourceID, status,
	)
	details := []string{}
	if data.Status != "" && data.Status != "healthy" {
		details = append(details, "status: "+formatContinuousBackupStatus(data.Status))
	}
	if data.Error == "" {
		if len(details) == 0 {
			return message
		}

		return fmt.Sprintf("%s (%s)", message, strings.Join(details, ", "))
	}

	details = append(details, "error: "+data.Error)
	return fmt.Sprintf("%s (%s)", message, strings.Join(details, ", "))
}

func formatContinuousBackupStatus(status string) string {
	switch status {
	case "pgbackrest_error":
		return "pgBackRest error"
	case "wal_error":
		return "WAL error"
	default:
		return status
	}
}
