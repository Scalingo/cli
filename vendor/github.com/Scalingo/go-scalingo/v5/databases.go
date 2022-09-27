package scalingo

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v5/http"
)

type DatabasesService interface {
	DatabaseShow(ctx context.Context, app, addonID string) (Database, error)
	DatabaseEnableFeature(ctx context.Context, app, addonID, feature string) (DatabaseEnableFeatureResponse, error)
	DatabaseDisableFeature(ctx context.Context, app, addonID, feature string) (DatabaseDisableFeatureResponse, error)
}

type DatabaseStatus string

const (
	DatabaseStatusCreating  DatabaseStatus = "creating"
	DatabaseStatusRunning   DatabaseStatus = "running"
	DatabaseStatusMigrating DatabaseStatus = "migrating"
	DatabaseStatusUpdating  DatabaseStatus = "updating"
	DatabaseStatusUpgrading DatabaseStatus = "upgrading"
	DatabaseStatusStopped   DatabaseStatus = "stopped"
)

type DatabaseFeature struct {
	Name   string                `json:"name"`
	Status DatabaseFeatureStatus `json:"status"`
}

type Database struct {
	ID                         string            `json:"id"`
	CreatedAt                  time.Time         `json:"created_at"`
	ResourceID                 string            `json:"resource_id"`
	AppName                    string            `json:"app_name"`
	EncryptionAtRest           bool              `json:"encryption_at_rest"`
	Features                   []DatabaseFeature `json:"features"`
	Plan                       string            `json:"plan"`
	Status                     DatabaseStatus    `json:"status"`
	TypeID                     string            `json:"type_id"`
	TypeName                   string            `json:"type_name"`
	VersionID                  string            `json:"version_id"`
	MongoReplSetName           string            `json:"mongo_repl_set_name"`
	Instances                  []Instance        `json:"instances"`
	NextVersionID              string            `json:"next_version_id"`
	ReadableVersion            string            `json:"readable_version"`
	Hostname                   string            `json:"hostname"`
	CurrentOperationID         string            `json:"current_operation_id"`
	Cluster                    bool              `json:"cluster"`
	PeriodicBackupsEnabled     bool              `json:"periodic_backups_enabled"`
	PeriodicBackupsScheduledAt []int             `json:"periodic_backups_scheduled_at"` // Hours of the day of the periodic backups (UTC)
}

type InstanceStatus string

const (
	InstanceStatusBooting    InstanceStatus = "booting"
	InstanceStatusRunning    InstanceStatus = "running"
	InstanceStatusRestarting InstanceStatus = "restarting"
	InstanceStatusMigrating  InstanceStatus = "migrating"
	InstanceStatusUpgrading  InstanceStatus = "upgrading"
	InstanceStatusStopped    InstanceStatus = "stopped"
	InstanceStatusRemoving   InstanceStatus = "removing"
)

type InstanceType string

const (
	// InstanceTypeDBNode instances represent database instances holding data
	InstanceTypeDBNode InstanceType = "db-node"
	// InstanceTypeUtility instances are those running service used for running a
	// service which is neither holding data nor routing requests utilities as
	// stated by its Name
	InstanceTypeUtility InstanceType = "utility"
	// InstanceTypeHAProxy are instances running a HAProxy reverse proxy in order
	// to route requests to the InstanceTypeDBNodes
	InstanceTypeHAProxy InstanceType = "haproxy"
)

type Instance struct {
	ID        string         `json:"id"`
	Hostname  string         `json:"hostname"`
	Port      int            `json:"port"`
	Status    InstanceStatus `json:"status"`
	Type      InstanceType   `json:"type"`
	Features  []string       `json:"features"`
	PrivateIP string         `json:"private_ip"`
}

type DatabaseRes struct {
	Database Database `json:"database"`
}

func (c *Client) DatabaseShow(ctx context.Context, app, addonID string) (Database, error) {
	var res DatabaseRes
	err := c.DBAPI(app, addonID).ResourceGet(ctx, "databases", addonID, nil, &res)
	if err != nil {
		return Database{}, errgo.Notef(err, "fail to get the database")
	}
	return res.Database, nil
}

type PeriodicBackupsConfigParams struct {
	ScheduledAt *int  `json:"periodic_backups_scheduled_at,omitempty"`
	Enabled     *bool `json:"periodic_backups_enabled,omitempty"`
}

func (c *Client) PeriodicBackupsConfig(ctx context.Context, app, addonID string, params PeriodicBackupsConfigParams) (Database, error) {
	var dbRes DatabaseRes
	err := c.DBAPI(app, addonID).ResourceUpdate(ctx, "databases", addonID, map[string]PeriodicBackupsConfigParams{
		"database": params,
	}, &dbRes)
	if err != nil {
		return Database{}, errgo.Notef(err, "fail to update periodic backups settings")
	}
	return dbRes.Database, nil
}

type DatabaseEnableFeatureParams struct {
	Feature DatabaseFeature `json:"feature"`
}

type DatabaseFeatureStatus string

const (
	DatabaseFeatureStatusActivated DatabaseFeatureStatus = "ACTIVATED"
	DatabaseFeatureStatusPending   DatabaseFeatureStatus = "PENDING"
	DatabaseFeatureStatusFailed    DatabaseFeatureStatus = "FAILED"
)

type DatabaseEnableFeatureResponse struct {
	Name    string                `json:"name"`
	Status  DatabaseFeatureStatus `json:"status"`
	Message string                `json:"message"`
}

func (c *Client) DatabaseEnableFeature(ctx context.Context, app, addonID, feature string) (DatabaseEnableFeatureResponse, error) {
	payload := DatabaseEnableFeatureParams{
		Feature: DatabaseFeature{
			Name: feature,
		},
	}

	res := DatabaseEnableFeatureResponse{}
	err := c.DBAPI(app, addonID).DoRequest(ctx, &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/databases/" + addonID + "/features",
		Expected: httpclient.Statuses{http.StatusOK},
		Params:   payload,
	}, &res)

	if err != nil {
		return res, errgo.Notef(err, "fail to enable database feature %v", feature)
	}

	return res, nil
}

type DatabaseDisableFeatureResponse struct {
	Message string `json:"message"`
}

func (c *Client) DatabaseDisableFeature(ctx context.Context, app, addonID, feature string) (DatabaseDisableFeatureResponse, error) {
	res := DatabaseDisableFeatureResponse{}
	err := c.DBAPI(app, addonID).DoRequest(ctx, &httpclient.APIRequest{
		Method:   "DELETE",
		Endpoint: "/databases/" + addonID + "/features",
		Expected: httpclient.Statuses{http.StatusOK},
		Params:   map[string]string{"feature": feature},
	}, &res)

	if err != nil {
		return res, errgo.Notef(err, "fail to disable database feature %v", feature)
	}

	return res, nil
}
