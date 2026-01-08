package scalingo

import (
	"fmt"
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
	Feature           string `json:"feature"`
	AddonProviderID   string `json:"addon_provider_id"`
	AddonProviderName string `json:"addon_provider_name"`
	AddonUUID         string `json:"addon_uuid"`
	EventSecurityTypeData
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
	Feature           string `json:"feature"`
	AddonProviderID   string `json:"addon_provider_id"`
	AddonProviderName string `json:"addon_provider_name"`
	AddonUUID         string `json:"addon_uuid"`
	EventSecurityTypeData
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
	AddonUUID    string       `json:"addon_uuid"`
	AddonName    string       `json:"addon_name"`
	ResourceID   string       `json:"resource_id"`
	BackupMethod BackupMethod `json:"backup_method"`
	BackupID     string       `json:"backup_id"`
	BackupStatus string       `json:"backup_status"`
	StartedAt    time.Time    `json:"started_at"`
	EndedAt      time.Time    `json:"ended_at"`
	EventSecurityTypeData
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
