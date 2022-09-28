package scalingo

import (
	"time"
)

type DataAccessConsent struct {
	AppID           string     `json:"app_id"`
	UserID          string     `json:"user_id"`
	ContainersUntil *time.Time `json:"containers_until,omitempty"`
	DatabasesUntil  *time.Time `json:"databases_until,omitempty"`
}
