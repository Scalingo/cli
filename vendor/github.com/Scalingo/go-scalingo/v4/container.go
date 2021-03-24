package scalingo

import (
	"time"
)

type Container struct {
	ID        string     `json:"id"`
	AppID     string     `json:"app_id"`
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Command   string     `json:"command"`
	Type      string     `json:"type"`
	TypeIndex int        `json:"type_index"`
	Label     string     `json:"label"`
	State     string     `json:"state"`
	Size      string     `json:"size"`
	App       *App       `json:"app"`
}
