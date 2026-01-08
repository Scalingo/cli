package scalingo

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/errors/v2"
)

type MaintenanceWindow struct {
	WeekdayUTC      int `json:"weekday_utc"`
	StartingHourUTC int `json:"starting_hour_utc"`
	DurationInHour  int `json:"duration_in_hour"`
}

type MaintenanceWindowParams struct {
	WeekdayUTC      *int `json:"weekday_utc,omitempty"`
	StartingHourUTC *int `json:"starting_hour_utc,omitempty"`
}

type Maintenance struct {
	ID         string            `json:"id"`
	DatabaseID string            `json:"database_id"`
	Status     MaintenanceStatus `json:"status"`
	Type       string            `json:"type"`
	StartedAt  *time.Time        `json:"started_at,omitempty"`
	EndedAt    *time.Time        `json:"ended_at,omitempty"`
}

type MaintenanceStatus string

const (
	MaintenanceStatusScheduled MaintenanceStatus = "scheduled"
	MaintenanceStatusNotified  MaintenanceStatus = "notified"
	MaintenanceStatusQueued    MaintenanceStatus = "queued"
	MaintenanceStatusCancelled MaintenanceStatus = "cancelled"
	MaintenanceStatusRunning   MaintenanceStatus = "running"
	MaintenanceStatusFailed    MaintenanceStatus = "failed"
	MaintenanceStatusDone      MaintenanceStatus = "done"
)

func (c *Client) DatabaseUpdateMaintenanceWindow(ctx context.Context, app, addonID string, params MaintenanceWindowParams) (Database, error) {
	var dbRes DatabaseRes
	err := c.DBAPI(app, addonID).ResourceUpdate(ctx, "databases", addonID, map[string]interface{}{
		"database": map[string]interface{}{
			"maintenance_window": params,
		},
	}, &dbRes)

	if err != nil {
		return Database{}, errors.Notef(ctx, err, "update database '%v' maintenance window", addonID)
	}
	return dbRes.Database, nil
}

type ListMaintenanceResponse struct {
	Maintenance []*Maintenance `json:"maintenance"`
	Meta        struct {
		PaginationMeta PaginationMeta `json:"pagination"`
	}
}

func (c *Client) DatabaseListMaintenance(ctx context.Context, app, addonID string, opts PaginationOpts) ([]*Maintenance, PaginationMeta, error) {
	var maintenanceRes ListMaintenanceResponse
	err := c.DBAPI(app, addonID).SubresourceList(ctx, "databases", addonID, "maintenance", opts.ToMap(), &maintenanceRes)
	if err != nil {
		return nil, PaginationMeta{}, errors.Notef(ctx, err, "list database '%v' maintenance", addonID)
	}
	return maintenanceRes.Maintenance, maintenanceRes.Meta.PaginationMeta, nil
}

func (c *Client) DatabaseShowMaintenance(ctx context.Context, app, addonID, maintenanceID string) (Maintenance, error) {
	var maintenanceRes Maintenance
	err := c.DBAPI(app, addonID).SubresourceGet(ctx, "databases", addonID, "maintenance", maintenanceID, nil, &maintenanceRes)
	if err != nil {
		return maintenanceRes, errors.Notef(ctx, err, "get database '%v' maintenance", addonID)
	}
	return maintenanceRes, nil
}
