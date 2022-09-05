package scalingo

import "fmt"

type EventLoginSuccessType struct {
	Event
	TypeData EventLoginSuccessTypeData `json:"type_data"`
}
type EventLoginSuccessTypeData EventSecurityTypeData

func (ev *EventLoginSuccessType) String() string {
	return fmt.Sprintf("Successful login from %v", ev.TypeData.RemoteIP)
}

type EventLoginFailureType struct {
	Event
	TypeData EventLoginFailureTypeData `json:"type_data"`
}
type EventLoginFailureTypeData EventSecurityTypeData

func (ev *EventLoginFailureType) String() string {
	return fmt.Sprintf("Failed login attempt from %v", ev.TypeData.RemoteIP)
}

type EventLoginLockType struct {
	Event
	TypeData EventLoginLockTypeData `json:"type_data"`
}
type EventLoginLockTypeData EventSecurityTypeData

func (ev *EventLoginLockType) String() string {
	return "Account is locked"
}

type EventLoginUnlockSuccessType struct {
	Event
	TypeData EventLoginUnlockSuccessTypeData `json:"type_data"`
}
type EventLoginUnlockSuccessTypeData EventSecurityTypeData

func (ev *EventLoginUnlockSuccessType) String() string {
	return fmt.Sprintf("Account unlocked from %v", ev.TypeData.RemoteIP)
}

type EventPasswordResetQueryType struct {
	Event
	TypeData EventPasswordResetQueryTypeData `json:"type_data"`
}
type EventPasswordResetQueryTypeData EventSecurityTypeData

func (ev *EventPasswordResetQueryType) String() string {
	return fmt.Sprintf("Password reset process initiated from %v", ev.TypeData.RemoteIP)
}

type EventPasswordResetSuccessType struct {
	Event
	TypeData EventPasswordResetSuccessTypeData `json:"type_data"`
}
type EventPasswordResetSuccessTypeData EventSecurityTypeData

func (ev *EventPasswordResetSuccessType) String() string {
	return fmt.Sprintf("Password changed from %v", ev.TypeData.RemoteIP)
}
