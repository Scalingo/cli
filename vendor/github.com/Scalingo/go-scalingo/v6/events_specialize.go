package scalingo

// Do not edit, generated with 'go generate'

import (
	"encoding/json"

	"github.com/Scalingo/go-scalingo/v6/debug"
)

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
	case EventScale:
		e = &EventScaleType{Event: ev}
	case EventStopApp:
		e = &EventStopAppType{Event: ev}
	case EventCrash:
		e = &EventCrashType{Event: ev}
	case EventRepeatedCrash:
		e = &EventRepeatedCrashType{Event: ev}
	case EventDeployment:
		e = &EventDeploymentType{Event: ev}
	case EventLinkSCM:
		e = &EventLinkSCMType{Event: ev}
	case EventUpdateSCM:
		e = &EventUpdateSCMType{Event: ev}
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
	case EventRun:
		e = &EventRunType{Event: ev}
	case EventNewDomain:
		e = &EventNewDomainType{Event: ev}
	case EventEditDomain:
		e = &EventEditDomainType{Event: ev}
	case EventDeleteDomain:
		e = &EventDeleteDomainType{Event: ev}
	case EventUpgradeDatabase:
		e = &EventUpgradeDatabaseType{Event: ev}
	case EventNewAddon:
		e = &EventNewAddonType{Event: ev}
	case EventUpgradeAddon:
		e = &EventUpgradeAddonType{Event: ev}
	case EventDeleteAddon:
		e = &EventDeleteAddonType{Event: ev}
	case EventResumeAddon:
		e = &EventResumeAddonType{Event: ev}
	case EventSuspendAddon:
		e = &EventSuspendAddonType{Event: ev}
	case EventDatabaseAddFeature:
		e = &EventDatabaseAddFeatureType{Event: ev}
	case EventDatabaseRemoveFeature:
		e = &EventDatabaseRemoveFeatureType{Event: ev}
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
	case EventEditKey:
		e = &EventEditKeyType{Event: ev}
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
	case EventEditAutoscaler:
		e = &EventEditAutoscalerType{Event: ev}
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
	case EventNewNotifier:
		e = &EventNewNotifierType{Event: ev}
	case EventEditNotifier:
		e = &EventEditNotifierType{Event: ev}
	case EventDeleteNotifier:
		e = &EventDeleteNotifierType{Event: ev}
	case EventEditHDSContact:
		e = &EventEditHDSContactType{Event: ev}
	case EventCreateDataAccessConsent:
		e = &EventCreateDataAccessConsentType{Event: ev}
	case EventNewToken:
		e = &EventNewTokenType{Event: ev}
	case EventRegenerateToken:
		e = &EventRegenerateTokenType{Event: ev}
	case EventDeleteToken:
		e = &EventDeleteTokenType{Event: ev}
	case EventTfaEnabled:
		e = &EventTfaEnabledType{Event: ev}
	case EventTfaDisabled:
		e = &EventTfaDisabledType{Event: ev}
	case EventLoginSuccess:
		e = &EventLoginSuccessType{Event: ev}
	case EventLoginFailure:
		e = &EventLoginFailureType{Event: ev}
	case EventLoginLock:
		e = &EventLoginLockType{Event: ev}
	case EventLoginUnlockSuccess:
		e = &EventLoginUnlockSuccessType{Event: ev}
	case EventPasswordResetQuery:
		e = &EventPasswordResetQueryType{Event: ev}
	case EventPasswordResetSuccess:
		e = &EventPasswordResetSuccessType{Event: ev}
	case EventStackChanged:
		e = &EventStackChangedType{Event: ev}
	case EventCreateReviewApp:
		e = &EventCreateReviewAppType{Event: ev}
	case EventDestroyReviewApp:
		e = &EventDestroyReviewAppType{Event: ev}
	case EventPlanDatabaseMaintenance:
		e = &EventPlanDatabaseMaintenanceType{Event: ev}
	case EventStartDatabaseMaintenance:
		e = &EventStartDatabaseMaintenanceType{Event: ev}
	case EventCompleteDatabaseMaintenance:
		e = &EventCompleteDatabaseMaintenanceType{Event: ev}
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
