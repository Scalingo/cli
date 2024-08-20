package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v7/http"
)

const (
	RegionMigrationStatusCreated          RegionMigrationStatus = "created"
	RegionMigrationStatusPreflightSuccess RegionMigrationStatus = "preflight-success"
	RegionMigrationStatusPreflightError   RegionMigrationStatus = "preflight-error"
	RegionMigrationStatusRunning          RegionMigrationStatus = "running"
	RegionMigrationStatusPrepared         RegionMigrationStatus = "prepared"
	RegionMigrationStatusDataMigrated     RegionMigrationStatus = "data-migrated"
	RegionMigrationStatusAborting         RegionMigrationStatus = "aborting"
	RegionMigrationStatusAborted          RegionMigrationStatus = "aborted"
	RegionMigrationStatusError            RegionMigrationStatus = "error"
	RegionMigrationStatusDone             RegionMigrationStatus = "done"

	RegionMigrationStepAbort     RegionMigrationStep = "abort"
	RegionMigrationStepPreflight RegionMigrationStep = "preflight"
	RegionMigrationStepPrepare   RegionMigrationStep = "prepare"
	RegionMigrationStepData      RegionMigrationStep = "data"
	RegionMigrationStepFinalize  RegionMigrationStep = "finalize"

	StepStatusRunning StepStatus = "running"
	StepStatusDone    StepStatus = "done"
	StepStatusError   StepStatus = "error"
)

type RegionMigrationsService interface {
	CreateRegionMigration(ctx context.Context, appID string, params RegionMigrationParams) (RegionMigration, error)
	RunRegionMigrationStep(ctx context.Context, appID, migrationID string, step RegionMigrationStep) error
	ShowRegionMigration(ctx context.Context, appID, migrationID string) (RegionMigration, error)
	ListRegionMigrations(ctx context.Context, appID string) ([]RegionMigration, error)
}

type RegionMigrationParams struct {
	Destination string `json:"destination"`
	DstAppName  string `json:"dst_app_name"`
}

type RegionMigration struct {
	ID          string                `json:"id"`
	SrcAppName  string                `json:"src_app_name"`
	DstAppName  string                `json:"dst_app_name"`
	AppID       string                `json:"app_id"`
	NewAppID    string                `json:"new_app_id"`
	Source      string                `json:"source"`
	Destination string                `json:"destination"`
	Status      RegionMigrationStatus `json:"status"`
	StartedAt   time.Time             `json:"started_at"`
	FinishedAt  time.Time             `json:"finished_at"`
	Steps       Steps                 `json:"steps"`
}

type StepStatus string
type RegionMigrationStatus string
type RegionMigrationStep string
type Steps []Step

type Step struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Status StepStatus `json:"status"`
	Logs   string     `json:"logs"`
}

func (c *Client) CreateRegionMigration(ctx context.Context, appID string, params RegionMigrationParams) (RegionMigration, error) {
	var migration RegionMigration

	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", appID, "region_migrations", map[string]RegionMigrationParams{
		"migration": params,
	}, &migration)
	if err != nil {
		return migration, errgo.Notef(err, "fail to create migration")
	}

	return migration, nil
}

func (c *Client) RunRegionMigrationStep(ctx context.Context, appID, migrationID string, step RegionMigrationStep) error {
	err := c.ScalingoAPI().DoRequest(ctx, &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + appID + "/region_migrations/" + migrationID + "/run",
		Params:   map[string]RegionMigrationStep{"step": step},
		Expected: http.Statuses{204},
	}, nil)
	if err != nil {
		return errgo.Notef(err, "fail to run migration step")
	}
	return nil
}

func (c *Client) ShowRegionMigration(ctx context.Context, appID, migrationID string) (RegionMigration, error) {
	var migration RegionMigration

	err := c.ScalingoAPI().SubresourceGet(ctx, "apps", appID, "region_migrations", migrationID, nil, &migration)
	if err != nil {
		return migration, errgo.Notef(err, "fail to get migration")
	}

	return migration, nil
}

func (c *Client) ListRegionMigrations(ctx context.Context, appID string) ([]RegionMigration, error) {
	var migrations []RegionMigration

	err := c.ScalingoAPI().SubresourceList(ctx, "apps", appID, "region_migrations", nil, &migrations)
	if err != nil {
		return migrations, errgo.Notef(err, "fail to list migrations")
	}

	return migrations, nil
}
