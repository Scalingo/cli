package scalingo

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Scalingo/go-scalingo/debug"
)

type Event struct {
	ID          string                 `json:"id"`
	AppID       string                 `json:"app_id"`
	CreatedAt   time.Time              `json:"created_at"`
	User        EventUser              `json:"user"`
	Type        EventTypeName          `json:"type"`
	AppName     string                 `json:"app_name"`
	RawTypeData json.RawMessage        `json:"type_data"`
	TypeData    map[string]interface{} `json:"-"`
}

func (ev *Event) GetEvent() *Event {
	return ev
}

func (ev *Event) TypeDataPtr() interface{} {
	return ev.TypeData
}

func (ev *Event) String() string {
	return fmt.Sprintf("Unknown event %v on app %v", ev.Type, ev.AppName)
}

func (ev *Event) When() string {
	return ev.CreatedAt.Format("Mon Jan 02 2006 15:04:05")
}

func (ev *Event) Who() string {
	return fmt.Sprintf("%s (%s)", ev.User.Username, ev.User.Email)
}

func (ev *Event) PrintableType() string {
	return strings.Title(strings.Replace(string(ev.Type), "_", " ", -1))
}

type DetailedEvent interface {
	fmt.Stringer
	GetEvent() *Event
	PrintableType() string
	When() string
	Who() string
	TypeDataPtr() interface{}
}

type Events []DetailedEvent

type EventUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

type EventTypeName string

const (
	EventNewUser              EventTypeName = "new_user"
	EventNewApp               EventTypeName = "new_app"
	EventEditApp              EventTypeName = "edit_app"
	EventDeleteApp            EventTypeName = "delete_app"
	EventRenameApp            EventTypeName = "rename_app"
	EventTransferApp          EventTypeName = "transfer_app"
	EventRestart              EventTypeName = "restart"
	EventScale                EventTypeName = "scale"
	EventStopApp              EventTypeName = "stop_app"
	EventCrash                EventTypeName = "crash"
	EventDeployment           EventTypeName = "deployment"
	EventLinkSCM              EventTypeName = "link_scm"
	EventUnlinkSCM            EventTypeName = "unlink_scm"
	EventNewIntegration       EventTypeName = "new_integration"
	EventDeleteIntegration    EventTypeName = "delete_integration"
	EventAuthorizeGithub      EventTypeName = "authorize_github"
	EventRevokeGithub         EventTypeName = "revoke_github"
	EventRun                  EventTypeName = "run"
	EventNewDomain            EventTypeName = "new_domain"
	EventEditDomain           EventTypeName = "edit_domain"
	EventDeleteDomain         EventTypeName = "delete_domain"
	EventNewAddon             EventTypeName = "new_addon"
	EventUpgradeAddon         EventTypeName = "upgrade_addon"
	EventUpgradeDatabase      EventTypeName = "upgrade_database"
	EventDeleteAddon          EventTypeName = "delete_addon"
	EventResumeAddon          EventTypeName = "resume_addon"
	EventSuspendAddon         EventTypeName = "suspend_addon"
	EventNewCollaborator      EventTypeName = "new_collaborator"
	EventAcceptCollaborator   EventTypeName = "accept_collaborator"
	EventDeleteCollaborator   EventTypeName = "delete_collaborator"
	EventNewVariable          EventTypeName = "new_variable"
	EventEditVariable         EventTypeName = "edit_variable"
	EventEditVariables        EventTypeName = "edit_variables"
	EventDeleteVariable       EventTypeName = "delete_variable"
	EventAddCredit            EventTypeName = "add_credit"
	EventAddPaymentMethod     EventTypeName = "add_payment_method"
	EventAddVoucher           EventTypeName = "add_voucher"
	EventNewKey               EventTypeName = "new_key"
	EventDeleteKey            EventTypeName = "delete_key"
	EventPaymentAttempt       EventTypeName = "payment_attempt"
	EventNewAlert             EventTypeName = "new_alert"
	EventAlert                EventTypeName = "alert"
	EventDeleteAlert          EventTypeName = "delete_alert"
	EventNewAutoscaler        EventTypeName = "new_autoscaler"
	EventDeleteAutoscaler     EventTypeName = "delete_autoscaler"
	EventAddonUpdated         EventTypeName = "addon_updated"
	EventStartRegionMigration EventTypeName = "start_region_migration"
	EventNewLogDrain          EventTypeName = "new_log_drain"
	EventDeleteLogDrain       EventTypeName = "delete_log_drain"
	EventNewAddonLogDrain     EventTypeName = "new_addon_log_drain"
	EventDeleteAddonLogDrain  EventTypeName = "delete_addon_log_drain"

	// EventLinkGithub and EventUnlinkGithub events are kept for
	// retro-compatibility. They are replaced by SCM events.
	EventLinkGithub   EventTypeName = "link_github"
	EventUnlinkGithub EventTypeName = "unlink_github"
)

type EventNewUserType struct {
	Event
	TypeData EventNewUserTypeData `json:"type_data"`
}

func (ev *EventNewUserType) String() string {
	return fmt.Sprintf("You joined Scalingo. Hooray!")
}

type EventNewUserTypeData struct {
}

type EventNewAppType struct {
	Event
	TypeData EventNewAppTypeData `json:"type_data"`
}

func (ev *EventNewAppType) String() string {
	return fmt.Sprintf("the application has been created")
}

type EventNewAppTypeData struct {
	GitSource string `json:"git_source"`
}

type EventEditAppType struct {
	Event
	TypeData EventEditAppTypeData `json:"type_data"`
}

type EventEditAppTypeData struct {
	ForceHTTPS *bool `json:"force_https"`
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
	return fmt.Sprintf("the application has been deleted")
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

type EventRenameAppTypeData struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type EventTransferAppType struct {
	Event
	TypeData EventTransferAppTypeData `json:"type_data"`
}

func (ev *EventTransferAppType) String() string {
	return fmt.Sprintf(
		"the application has been transfered to %s (%s)",
		ev.TypeData.NewOwner.Username, ev.TypeData.NewOwner.Email,
	)
}

type EventTransferAppTypeData struct {
	OldOwner EventUser `json:"old_owner"`
	NewOwner EventUser `json:"new_owner"`
}

type EventRestartType struct {
	Event
	TypeData EventRestartTypeData `json:"type_data"`
}

func (ev *EventRestartType) String() string {
	if len(ev.TypeData.Scope) != 0 {
		return fmt.Sprintf("containers %v have been restarted", ev.TypeData.Scope)
	}
	return fmt.Sprintf("containers have been restarted")
}

func (ev *EventRestartType) Who() string {
	if ev.TypeData.AddonName != "" {
		return fmt.Sprintf("Addon %s", ev.TypeData.AddonName)
	} else {
		return ev.Event.Who()
	}
}

type EventRestartTypeData struct {
	Scope     []string `json:"scope"`
	AddonName string   `json:"addon_name"`
}

type EventStopAppType struct {
	Event
	TypeData EventStopAppTypeData `json:"type_data"`
}

func (ev *EventStopAppType) String() string {
	return fmt.Sprintf("app has been stopped (reason: %s)", ev.TypeData.Reason)
}

type EventStopAppTypeData struct {
	Reason string `json:"reason"`
}

func (ev *EventScaleType) String() string {
	return fmt.Sprintf(
		"containers have been scaled from %s, to %s",
		ev.TypeData.containersString(ev.TypeData.PreviousContainers),
		ev.TypeData.containersString(ev.TypeData.Containers),
	)
}

type EventScaleType struct {
	Event
	TypeData EventScaleTypeData `json:"type_data"`
}

type EventScaleTypeData struct {
	PreviousContainers map[string]string `json:"previous_containers"`
	Containers         map[string]string `json:"containers"`
}

func (e *EventScaleTypeData) containersString(containers map[string]string) string {
	types := []string{}
	for name, amountAndSize := range containers {
		types = append(types, fmt.Sprintf("%s:%s", name, amountAndSize))
	}
	return "[" + strings.Join(types, ", ") + "]"
}

type EventCrashType struct {
	Event
	TypeData EventCrashTypeData `json:"type_data"`
}

func (ev *EventCrashType) String() string {
	msg := fmt.Sprintf("container '%v' has crashed", ev.TypeData.ContainerType)

	if ev.TypeData.CrashLogs != "" {
		msg += fmt.Sprintf(" (logs on %s)", ev.TypeData.LogsUrl)
	}

	return msg
}

type EventCrashTypeData struct {
	ContainerType string `json:"container_type"`
	CrashLogs     string `json:"crash_logs"`
	LogsUrl       string `json:"logs_url"`
}

type EventDeploymentType struct {
	Event
	TypeData EventDeploymentTypeData `json:"type_data"`
}

func (ev *EventDeploymentType) String() string {
	return fmt.Sprintf("deployment of %s (%s)", ev.TypeData.GitRef, ev.TypeData.Status)
}

type EventDeploymentTypeData struct {
	DeploymentID string `json:"deployment_id"`
	Puser        string `json:"pusher"`
	GitRef       string `json:"git_ref"`
	Status       string `json:"status"`
	Duration     int    `json:"duration"`
}

type EventLinkGithubType struct {
	Event
	TypeData EventLinkGithubTypeData `json:"type_data"`
}

func (ev *EventLinkGithubType) String() string {
	return fmt.Sprintf("app has been linked to Github repository '%s'", ev.TypeData.RepoName)
}

type EventLinkGithubTypeData struct {
	RepoName       string `json:"repo_name"`
	LinkerUsername string `json:"linker_username"`
	GithubSource   string `json:"github_source"`
}

type EventUnlinkGithubType struct {
	Event
	TypeData EventUnlinkGithubTypeData `json:"type_data"`
}

func (ev *EventUnlinkGithubType) String() string {
	return fmt.Sprintf("app has been unlinked from Github repository '%s'", ev.TypeData.RepoName)
}

type EventUnlinkGithubTypeData struct {
	RepoName         string `json:"repo_name"`
	UnlinkerUsername string `json:"unlinker_username"`
	GithubSource     string `json:"github_source"`
}

type EventLinkSCMType struct {
	Event
	TypeData EventLinkSCMTypeData `json:"type_data"`
}

func (ev *EventLinkSCMType) String() string {
	return fmt.Sprintf("app has been linked to repository '%s'", ev.TypeData.RepoName)
}

type EventLinkSCMTypeData struct {
	RepoName       string `json:"repo_name"`
	LinkerUsername string `json:"linker_username"`
	Source         string `json:"source"`
}

type EventUnlinkSCMType struct {
	Event
	TypeData EventUnlinkSCMTypeData `json:"type_data"`
}

func (ev *EventUnlinkSCMType) String() string {
	return fmt.Sprintf("app has been unlinked from repository '%s'", ev.TypeData.RepoName)
}

type EventUnlinkSCMTypeData struct {
	RepoName         string `json:"repo_name"`
	UnlinkerUsername string `json:"unlinker_username"`
	Source           string `json:"source"`
}

type EventRunType struct {
	Event
	TypeData EventRunTypeData `json:"type_data"`
}

func (ev *EventRunType) String() string {
	return fmt.Sprintf("one-off container with command '%s'", ev.TypeData.Command)
}

type EventRunTypeData struct {
	Command    string `json:"command"`
	AuditLogID string `json:"audit_log_id"`
}

type EventNewDomainType struct {
	Event
	TypeData EventNewDomainTypeData `json:"type_data"`
}

func (ev *EventNewDomainType) String() string {
	return fmt.Sprintf("'%s' has been associated", ev.TypeData.Hostname)
}

type EventNewDomainTypeData struct {
	Hostname string `json:"hostname"`
	SSL      bool   `json:"ssl"`
}

type EventEditDomainType struct {
	Event
	TypeData EventEditDomainTypeData `json:"type_data"`
}

func (ev *EventEditDomainType) String() string {
	t := ev.TypeData
	res := fmt.Sprintf("'%s' modified", t.Hostname)
	if !t.SSL && t.OldSSL {
		res += ", TLS certificate has been removed"
	} else if !t.SSL && t.OldSSL {
		res += ", TLS certificate has been added"
	} else if t.SSL && t.OldSSL {
		res += ", TLS certificate has been changed"
	}
	return res
}

type EventEditDomainTypeData struct {
	Hostname string `json:"hostname"`
	SSL      bool   `json:"ssl"`
	OldSSL   bool   `json:"old_ssl"`
}

type EventDeleteDomainType struct {
	Event
	TypeData EventDeleteDomainTypeData `json:"type_data"`
}

func (ev *EventDeleteDomainType) String() string {
	return fmt.Sprintf("'%s' has been disassociated", ev.TypeData.Hostname)
}

type EventDeleteDomainTypeData struct {
	Hostname string `json:"hostname"`
}

type EventAddon struct {
	AddonProviderName string `json:"addon_provider_name"`
	PlanName          string `json:"plan_name"`
	ResourceID        string `json:"resource_id"`
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

type EventNewAddonTypeData struct {
	EventAddon
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

type EventUpgradeAddonTypeData struct {
	EventAddon
	OldPlanName string `json:"old_plan_name"`
	NewPlanName string `json:"new_plan_name"`
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

type EventResumeAddonTypeData struct {
	EventAddon
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

type EventSuspendAddonTypeData struct {
	EventAddon
	Reason string `json:"reason"`
}

type EventDeleteAddonTypeData struct {
	EventAddon
}

type EventCollaborator struct {
	EventUser
	Inviter EventUser `json:"inviter"`
}

type EventNewCollaboratorType struct {
	Event
	TypeData EventNewCollaboratorTypeData `json:"type_data"`
}

func (ev *EventNewCollaboratorType) String() string {
	return fmt.Sprintf(
		"'%s' has been invited",
		ev.TypeData.Collaborator.Email,
	)
}

type EventNewCollaboratorTypeData struct {
	Collaborator EventCollaborator `json:"collaborator"`
}

type EventAcceptCollaboratorType struct {
	Event
	TypeData EventAcceptCollaboratorTypeData `json:"type_data"`
}

func (ev *EventAcceptCollaboratorType) String() string {
	return fmt.Sprintf(
		"'%s' (%s) has accepted the collaboration",
		ev.TypeData.Collaborator.Username,
		ev.TypeData.Collaborator.Email,
	)
}

// Inviter is filled there
type EventAcceptCollaboratorTypeData struct {
	Collaborator EventCollaborator `json:"collaborator"`
}

type EventDeleteCollaboratorType struct {
	Event
	TypeData EventDeleteCollaboratorTypeData `json:"type_data"`
}

func (ev *EventDeleteCollaboratorType) String() string {
	return fmt.Sprintf(
		"'%s' (%s) is not a collaborator anymore",
		ev.TypeData.Collaborator.Username,
		ev.TypeData.Collaborator.Email,
	)
}

type EventDeleteCollaboratorTypeData struct {
	Collaborator EventCollaborator `json:"collaborator"`
}

type EventUpgradeDatabaseType struct {
	Event
	TypeData EventUpgradeDatabaseTypeData `json:"type_data"`
}

type EventUpgradeDatabaseTypeData struct {
	AddonName  string `json:"addon_name"`
	OldVersion string `json:"old_version"`
	NewVersion string `json:"new_version"`
}

func (ev *EventUpgradeDatabaseType) String() string {
	return fmt.Sprintf(
		"'%s' upgraded from v%s to v%s",
		ev.TypeData.AddonName, ev.TypeData.OldVersion, ev.TypeData.NewVersion,
	)
}

func (ev *EventUpgradeDatabaseType) Who() string {
	if ev.TypeData.AddonName != "" {
		return fmt.Sprintf("Addon %s", ev.TypeData.AddonName)
	} else {
		return ev.Event.Who()
	}
}

type EventVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EventNewVariableType struct {
	Event
	TypeData EventNewVariableTypeData `json:"type_data"`
}

func (ev *EventNewVariableType) String() string {
	return fmt.Sprintf("'%s' added to the environment", ev.TypeData.Name)
}

func (ev *EventNewVariableType) Who() string {
	if ev.TypeData.AddonName != "" {
		return fmt.Sprintf("Addon %s", ev.TypeData.AddonName)
	} else {
		return ev.Event.Who()
	}
}

type EventNewVariableTypeData struct {
	AddonName string `json:"addon_name"`
	EventVariable
}

type EventVariables []EventVariable

func (evs EventVariables) Names() string {
	names := []string{}
	for _, e := range evs {
		names = append(names, e.Name)
	}
	return strings.Join(names, ", ")
}

type EventEditVariableType struct {
	Event
	TypeData EventEditVariableTypeData `json:"type_data"`
}

func (ev *EventEditVariableType) String() string {
	return fmt.Sprintf("'%s' modified", ev.TypeData.Name)
}

type EventEditVariableTypeData struct {
	EventVariable
	OldValue  string `json:"old_value"`
	AddonName string `json:"addon_name"`
}

type EventEditVariablesType struct {
	Event
	TypeData EventEditVariablesTypeData `json:"type_data"`
}

func (ev *EventEditVariablesType) String() string {
	res := "environment changes:"
	if len(ev.TypeData.NewVars) > 0 {
		res += fmt.Sprintf(" %s added", ev.TypeData.NewVars.Names())
	}
	if len(ev.TypeData.UpdatedVars) > 0 {
		res += fmt.Sprintf(" %s modified", ev.TypeData.UpdatedVars.Names())
	}
	if len(ev.TypeData.DeletedVars) > 0 {
		res += fmt.Sprintf(" %s removed", ev.TypeData.DeletedVars.Names())
	}
	return res
}

func (ev *EventEditVariableType) Who() string {
	if ev.TypeData.AddonName != "" {
		return fmt.Sprintf("Addon %s", ev.TypeData.AddonName)
	} else {
		return ev.Event.Who()
	}
}

type EventEditVariablesTypeData struct {
	NewVars     EventVariables `json:"new_vars"`
	UpdatedVars EventVariables `json:"updated_vars"`
	DeletedVars EventVariables `json:"deleted_vars"`
}

type EventDeleteVariableType struct {
	Event
	TypeData EventDeleteVariableTypeData `json:"type_data"`
}

func (ev *EventDeleteVariableType) String() string {
	return fmt.Sprintf("'%s' removed from environment", ev.TypeData.Name)
}

type EventDeleteVariableTypeData struct {
	EventVariable
}

type EventPaymentAttemptTypeData struct {
	Amount        float32 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"`
}

type EventPaymentAttemptType struct {
	Event
	TypeData EventPaymentAttemptTypeData `json:"type_data"`
}

func (ev *EventPaymentAttemptType) String() string {
	res := "Payment attempt of "
	res += fmt.Sprintf("%0.2f€", ev.TypeData.Amount)
	res += " with your "
	if ev.TypeData.PaymentMethod == "credit" {
		res += "credits"
	} else {
		res += "card"
	}
	if ev.TypeData.Status == "new" {
		res += " (pending)"
	} else if ev.TypeData.Status == "paid" {
		res += " (success)"
	} else {
		res += " (fail)"
	}
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

type EventNewAutoscalerTypeData struct {
	ContainerType string  `json:"container_type"`
	MinContainers int     `json:"min_containers"`
	MaxContainers int     `json:"max_containers"`
	Metric        string  `json:"metric"`
	Target        float64 `json:"target"`
	TargetText    string  `json:"target_text"`
}

type EventNewAutoscalerType struct {
	Event
	TypeData EventNewAutoscalerTypeData `json:"type_data"`
}

func (ev *EventNewAutoscalerType) String() string {
	d := ev.TypeData
	return fmt.Sprintf("Autoscaler created about %s on container %s (target: %s)", d.Metric, d.ContainerType, d.TargetText)
}

type EventDeleteAutoscalerTypeData struct {
	ContainerType string `json:"container_type"`
	Metric        string `json:"metric"`
}

type EventDeleteAutoscalerType struct {
	Event
	TypeData EventDeleteAutoscalerTypeData `json:"type_data"`
}

func (ev *EventDeleteAutoscalerType) String() string {
	d := ev.TypeData
	return fmt.Sprintf("Alert deleted about %s on container %s", d.Metric, d.ContainerType)
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

type EventStartRegionMigrationTypeData struct {
	MigrationID string `json:"migration_id"`
	Destination string `json:"destination"`
	Source      string `json:"source"`
	DstAppName  string `json:"dst_app_name"`
}

type EventStartRegionMigrationType struct {
	Event
	TypeData EventStartRegionMigrationTypeData `json:"type_data"`
}

func (ev *EventStartRegionMigrationType) String() string {
	return fmt.Sprintf("Application region migration started from %s to %s/%s", ev.TypeData.Source, ev.TypeData.Destination, ev.TypeData.DstAppName)
}

// New log drain
type EventNewLogDrainTypeData struct {
	URL string `json:"url"`
}

type EventNewLogDrainType struct {
	Event
	TypeData EventNewLogDrainTypeData `json:"type_data"`
}

func (ev *EventNewLogDrainType) String() string {
	return fmt.Sprintf("Log drain added on %s app", ev.AppName)
}

// Delete log drain
type EventDeleteLogDrainTypeData struct {
	URL string `json:"url"`
}

type EventDeleteLogDrainType struct {
	Event
	TypeData EventDeleteLogDrainTypeData `json:"type_data"`
}

func (ev *EventDeleteLogDrainType) String() string {
	return fmt.Sprintf("Log drain deleted on %s app", ev.AppName)
}

// New addon log drain
type EventNewAddonLogDrainTypeData struct {
	URL       string `json:"url"`
	AddonUUID string `json:"addon_uuid"`
	AddonName string `json:"addon_name"`
}

type EventNewAddonLogDrainType struct {
	Event
	TypeData EventNewAddonLogDrainTypeData `json:"type_data"`
}

func (ev *EventNewAddonLogDrainType) String() string {
	return fmt.Sprintf("Log drain added for %s addon on %s app", ev.TypeData.AddonName, ev.AppName)
}

// Delete addon log drain
type EventDeleteAddonLogDrainTypeData struct {
	URL       string `json:"url"`
	AddonUUID string `json:"addon_uuid"`
	AddonName string `json:"addon_name"`
}

type EventDeleteAddonLogDrainType struct {
	Event
	TypeData EventDeleteAddonLogDrainTypeData `json:"type_data"`
}

func (ev *EventDeleteAddonLogDrainType) String() string {
	return fmt.Sprintf("Log drain deleted on %s addon for %s app", ev.TypeData.AddonName, ev.AppName)
}

func (pev *Event) Specialize() DetailedEvent {
	var e DetailedEvent
	ev := *pev
	switch ev.Type {
	case EventNewUser:
		e = &EventNewUserType{Event: ev}
	case EventNewApp:
		e = &EventNewAppType{Event: ev}
	case EventEditApp:
		e = &EventEditAppType{Event: ev}
	case EventDeleteApp:
		e = &EventDeleteAppType{Event: ev}
	case EventRenameApp:
		e = &EventRenameAppType{Event: ev}
	case EventTransferApp:
		e = &EventTransferAppType{Event: ev}
	case EventRestart:
		e = &EventRestartType{Event: ev}
	case EventStopApp:
		e = &EventStopAppType{Event: ev}
	case EventScale:
		e = &EventScaleType{Event: ev}
	case EventCrash:
		e = &EventCrashType{Event: ev}
	case EventLinkSCM:
		e = &EventLinkSCMType{Event: ev}
	case EventUnlinkSCM:
		e = &EventUnlinkSCMType{Event: ev}
	case EventNewIntegration:
		e = &EventNewIntegrationType{Event: ev}
	case EventDeleteIntegration:
		e = &EventDeleteIntegrationType{Event: ev}
	case EventAuthorizeGithub:
		e = &EventAuthorizeGithubType{Event: ev}
	case EventRevokeGithub:
		e = &EventRevokeGithubType{Event: ev}
	case EventDeployment:
		e = &EventDeploymentType{Event: ev}
	case EventRun:
		e = &EventRunType{Event: ev}
	case EventNewDomain:
		e = &EventNewDomainType{Event: ev}
	case EventEditDomain:
		e = &EventEditDomainType{Event: ev}
	case EventDeleteDomain:
		e = &EventDeleteDomainType{Event: ev}
	case EventNewAddon:
		e = &EventNewAddonType{Event: ev}
	case EventUpgradeAddon:
		e = &EventUpgradeAddonType{Event: ev}
	case EventUpgradeDatabase:
		e = &EventUpgradeDatabaseType{Event: ev}
	case EventDeleteAddon:
		e = &EventDeleteAddonType{Event: ev}
	case EventResumeAddon:
		e = &EventResumeAddonType{Event: ev}
	case EventSuspendAddon:
		e = &EventSuspendAddonType{Event: ev}
	case EventNewCollaborator:
		e = &EventNewCollaboratorType{Event: ev}
	case EventAcceptCollaborator:
		e = &EventAcceptCollaboratorType{Event: ev}
	case EventDeleteCollaborator:
		e = &EventDeleteCollaboratorType{Event: ev}
	case EventNewVariable:
		e = &EventNewVariableType{Event: ev}
	case EventEditVariable:
		e = &EventEditVariableType{Event: ev}
	case EventEditVariables:
		e = &EventEditVariablesType{Event: ev}
	case EventDeleteVariable:
		e = &EventDeleteVariableType{Event: ev}
	case EventAddCredit:
		e = &EventAddCreditType{Event: ev}
	case EventAddPaymentMethod:
		e = &EventAddPaymentMethodType{Event: ev}
	case EventAddVoucher:
		e = &EventAddVoucherType{Event: ev}
	case EventNewKey:
		e = &EventNewKeyType{Event: ev}
	case EventDeleteKey:
		e = &EventDeleteKeyType{Event: ev}
	case EventPaymentAttempt:
		e = &EventPaymentAttemptType{Event: ev}
	case EventNewAlert:
		e = &EventNewAlertType{Event: ev}
	case EventAlert:
		e = &EventAlertType{Event: ev}
	case EventDeleteAlert:
		e = &EventDeleteAlertType{Event: ev}
	case EventNewAutoscaler:
		e = &EventNewAutoscalerType{Event: ev}
	case EventDeleteAutoscaler:
		e = &EventDeleteAutoscalerType{Event: ev}
	case EventAddonUpdated:
		e = &EventAddonUpdatedType{Event: ev}
	case EventStartRegionMigration:
		e = &EventStartRegionMigrationType{Event: ev}
	case EventNewLogDrain:
		e = &EventNewLogDrainType{Event: ev}
	case EventDeleteLogDrain:
		e = &EventDeleteLogDrainType{Event: ev}
	case EventNewAddonLogDrain:
		e = &EventNewAddonLogDrainType{Event: ev}
	case EventDeleteAddonLogDrain:
		e = &EventDeleteAddonLogDrainType{Event: ev}
	// Deprecated events. Replaced by equivalent with SCM in the name instead of
	// Github
	case EventLinkGithub:
		e = &EventLinkGithubType{Event: ev}
	case EventUnlinkGithub:
		e = &EventUnlinkGithubType{Event: ev}
	default:
		return pev
	}
	err := json.Unmarshal(pev.RawTypeData, e.TypeDataPtr())
	if err != nil {
		debug.Printf("error reading the data: %+v\n", err)
		return pev
	}
	return e
}
