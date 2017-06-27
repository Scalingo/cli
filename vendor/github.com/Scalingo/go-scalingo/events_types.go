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
	CreatedAt   time.Time              `json:"created_at"`
	User        EventUser              `json:"user"`
	Type        EventType              `json:"type"`
	AppID       string                 `json:"app_id"`
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
	return fmt.Sprintf("Unkown event %v on app %v", ev.Type, ev.AppName)
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

type EventType string

const (
	EventNewApp             EventType = "new_app"
	EventRenameApp                    = "rename_app"
	EventTransferApp                  = "transfer_app"
	EventRestart                      = "restart"
	EventScale                        = "scale"
	EventStopApp                      = "stop_app"
	EventCrash                        = "crash"
	EventDeployment                   = "deployment"
	EventLinkGithub                   = "link_github"
	EventUnlinkGithub                 = "unlink_github"
	EventRun                          = "run"
	EventNewDomain                    = "new_domain"
	EventEditDomain                   = "edit_domain"
	EventDeleteDomain                 = "delete_domain"
	EventNewAddon                     = "new_addon"
	EventUpgradeAddon                 = "upgrade_addon"
	EventUpgradeDatabase              = "upgrade_database"
	EventDeleteAddon                  = "delete_addon"
	EventResumeAddon                  = "resume_addon"
	EventSuspendAddon                 = "suspend_addon"
	EventNewCollaborator              = "new_collaborator"
	EventAcceptCollaborator           = "accept_collaborator"
	EventDeleteCollaborator           = "delete_collaborator"
	EventNewVariable                  = "new_variable"
	EventEditVariable                 = "edit_variable"
	EventEditVariables                = "edit_variables"
	EventDeleteVariable               = "delete_variable"
	EventNewNotification              = "new_notification"
	EventEditNotification             = "edit_notification"
	EventDeleteNotification           = "delete_notification"
	EventAddCredit                    = "add_credit"
	EventAddPaymentMethod             = "add_payment_method"
	EventAddVoucher                   = "add_voucher"
	EventAuthorizeGithub              = "authorize_github"
	EventRevokeGithub                 = "revoke_github"
	EventNewKey                       = "new_key"
	EventDeleteKey                    = "delete_key"
	EventPaymentAttempt               = "payment_attempt"
)

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
	} else {
		return fmt.Sprintf("containers have been restarted")
	}
}

type EventRestartTypeData struct {
	Scope         []string `json:"scope"`
	AddonProvider string   `json:"addon_provider"`
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
		dashboard_url := "https://my.scalingo.com/apps/" + ev.AppName + "/events/" + ev.ID
		msg += fmt.Sprintf(" (logs on %s)", dashboard_url)
	}

	return msg
}

type EventCrashTypeData struct {
	ContainerType string `json:"container_type"`
	CrashLogs     string `json:"crash_logs"`
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
	TypeData EventLinkGithubTypeData `json:"type_data"`
}

func (ev *EventUnlinkGithubType) String() string {
	return fmt.Sprintf("app has been unlinked from Github repository '%s'", ev.TypeData.RepoName)
}

type EventRunType struct {
	Event
	TypeData EventRunTypeData `json:"type_data"`
}

func (ev *EventRunType) String() string {
	return fmt.Sprintf("one-off container with command '%s'", ev.TypeData.Command)
}

type EventRunTypeData struct {
	Command string `json:"command"`
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
		return ev.Who()
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
		return ev.Who()
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
	OldValue string `json:"old_value"`
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

type EventNotification struct {
	NotificationType string `json:"notification_type"`
	Active           bool   `json:"active"`
	WebhookURL       string `json:"webhook_url"`
}

func (n *EventNotification) String() string {
	state := "disabled"
	if n.Active {
		state = "enabled"
	}
	return fmt.Sprintf("%s: %s (%s)", n.NotificationType, n.WebhookURL, state)
}

type EventNewNotificationType struct {
	Event
	TypeData EventNotification `json:"type_data"`
}

func (ev *EventNewNotificationType) String() string {
	return ev.TypeData.String()
}

type EventEditNotificationType struct {
	Event
	TypeData EventNotification `json:"type_data"`
}

func (ev *EventEditNotificationType) String() string {
	return ev.TypeData.String()
}

type EventDeleteNotificationType struct {
	Event
	TypeData EventNotification `json:"type_data"`
}

func (ev *EventDeleteNotificationType) String() string {
	return ev.TypeData.String()
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

func (pev *Event) Specialize() DetailedEvent {
	var e DetailedEvent
	ev := *pev
	switch ev.Type {
	case EventNewApp:
		e = &EventNewAppType{Event: ev}
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
	case EventLinkGithub:
		e = &EventLinkGithubType{Event: ev}
	case EventUnlinkGithub:
		e = &EventUnlinkGithubType{Event: ev}
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
	case EventNewNotification:
		e = &EventNewNotificationType{Event: ev}
	case EventEditNotification:
		e = &EventEditNotificationType{Event: ev}
	case EventDeleteNotification:
		e = &EventDeleteNotificationType{Event: ev}
	case EventAddCredit:
		e = &EventAddCreditType{Event: ev}
	case EventAddPaymentMethod:
		e = &EventAddPaymentMethodType{Event: ev}
	case EventAddVoucher:
		e = &EventAddVoucherType{Event: ev}
	case EventAuthorizeGithub:
		e = &EventAuthorizeGithubType{Event: ev}
	case EventNewKey:
		e = &EventNewKeyType{Event: ev}
	case EventDeleteKey:
		e = &EventDeleteKeyType{Event: ev}
	case EventPaymentAttempt:
		e = &EventPaymentAttemptType{Event: ev}
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
