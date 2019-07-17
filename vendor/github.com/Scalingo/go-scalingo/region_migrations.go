package scalingo

import (
	"time"

	errgo "gopkg.in/errgo.v1"
)

const (
	RegionMigrationStatusScheduled      RegionMigrationStatus = "scheduled"
	RegionMigrationStatusPreflightError RegionMigrationStatus = "preflight-error"
	RegionMigrationStatusRunning        RegionMigrationStatus = "running"
	RegionMigrationStatusError          RegionMigrationStatus = "error"
	RegionMigrationStatusDone           RegionMigrationStatus = "done"

	StepStatusRunning StepStatus = "running"
	StepStatusDone    StepStatus = "done"
	StepStatusError   StepStatus = "error"
)

type RegionMigrationsService interface {
	CreateRegionMigration(appID string, params RegionMigrationParams) (RegionMigration, error)
	ShowRegionMigration(appID, migrationID string) (RegionMigration, error)
	ListRegionMigrations(appID string) ([]RegionMigration, error)
}

type RegionMigrationParams struct {
	Destination string `json:"destination"`
}

type RegionMigration struct {
	ID          string                `json:"id"`
	AppName     string                `json:"app_name"`
	AppID       string                `json:"app_id"`
	NewAppID    string                `json:"new_app_id"`
	Destination string                `json:"destination"`
	Status      RegionMigrationStatus `json:"status"`
	StartedAt   time.Time             `json:"started_at"`
	FinishedAt  time.Time             `json:"finished_at"`
	Steps       Steps                 `json:"steps"`
}

type StepStatus string
type RegionMigrationStatus string
type Steps []Step

type Step struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Status StepStatus `json:"status"`
	Logs   string     `json:"logs"`
}

func (c *Client) CreateRegionMigration(appID string, params RegionMigrationParams) (RegionMigration, error) {
	var migration RegionMigration

	err := c.ScalingoAPI().SubresourceAdd("apps", appID, "region_migrations", map[string]RegionMigrationParams{
		"migration": params,
	}, &migration)
	if err != nil {
		return migration, errgo.Notef(err, "fail to create migration")
	}

	return migration, nil
}

func (c *Client) ShowRegionMigration(appID, migrationID string) (RegionMigration, error) {
	var migration RegionMigration

	err := c.ScalingoAPI().SubresourceGet("apps", appID, "region_migrations", migrationID, nil, &migration)
	if err != nil {
		return migration, errgo.Notef(err, "fail to get migration")
	}

	return migration, nil
}

func (c *Client) ListRegionMigrations(appID string) ([]RegionMigration, error) {
	var migrations []RegionMigration

	err := c.ScalingoAPI().SubresourceList("apps", appID, "region_migrations", nil, &migrations)
	if err != nil {
		return migrations, errgo.Notef(err, "fail to list migrations")
	}

	return migrations, nil
}
