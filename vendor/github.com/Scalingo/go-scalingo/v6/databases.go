package scalingo

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v6/http"
)

// DatabasesService is the interface gathering all the methods related to
// database addon configuration updates
type DatabasesService interface {
	DatabaseShow(ctx context.Context, app, addonID string) (Database, error)
	DatabaseEnableFeature(ctx context.Context, app, addonID, feature string) (DatabaseEnableFeatureResponse, error)
	DatabaseDisableFeature(ctx context.Context, app, addonID, feature string) (DatabaseDisableFeatureResponse, error)
	DatabaseUpdatePeriodicBackupsConfig(ctx context.Context, app, addonID string, params DatabaseUpdatePeriodicBackupsConfigParams) (Database, error)
}

// DatabaseStatus is a string representing the status of a database deployment
type DatabaseStatus string

const (
	// DatabaseStatusCreating is set when the database is being started before
	// it's operational
	DatabaseStatusCreating DatabaseStatus = "creating"
	// DatabaseStatusRunning is the standard status of a database when everything
	// is operational
	DatabaseStatusRunning DatabaseStatus = "running"
	// DatabaseStatusMigrating is set when a component of the database is being
	// migrated by Scalingo infrastructure
	DatabaseStatusMigrating DatabaseStatus = "migrating"
	// DatabaseStatusUpdating is set the plan of the database is being changed
	DatabaseStatusUpdating DatabaseStatus = "updating"
	// DatabaseStatusUpgrading is set when a database version upgrade is being
	// applied on the database
	DatabaseStatusUpgrading DatabaseStatus = "upgrading"
	// DatabaseStatusStopped is set when the database has been stopped (suspended
	// after free trial or when an account has been suspended)
	DatabaseStatusStopped DatabaseStatus = "stopped"
)

// DatabaseFeature represents the state of application of a database feature
type DatabaseFeature struct {
	Name   string                `json:"name"`
	Status DatabaseFeatureStatus `json:"status"`
}

// Database contains the metadata and configuration of a database deployment
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

// InstanceStatus is a type of string representing the status of an Instance
type InstanceStatus string

const (
	// InstanceStatusBooting is set when an instance is starting for the first
	// time
	InstanceStatusBooting InstanceStatus = "booting"
	// InstanceStatusRunning is the default status when the instance is working
	// normally
	InstanceStatusRunning InstanceStatus = "running"
	// InstanceStatusRestarting is set when an instance is restarting (during a
	// plan change, at the end of an upgrade or a migration)
	InstanceStatusRestarting InstanceStatus = "restarting"
	// InstanceStatusMigrating is set when an instance is being migrated by the
	// Scalingo infrastructure
	InstanceStatusMigrating InstanceStatus = "migrating"
	// InstanceStatusUpgrading is set when an instance version is being changed,
	// part of a Database upgrade
	InstanceStatusUpgrading InstanceStatus = "upgrading"
	// InstanceStatusStopped is set when an instance has been stopped
	InstanceStatusStopped InstanceStatus = "stopped"
	// InstanceStatusRemoving is set during the deletion of an Instance (business
	// to starter downgrade or database deletion)
	InstanceStatusRemoving InstanceStatus = "removing"
)

// InstanceType is a type of string representing the type of the Instance inside a Database
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

// Instance contains the metadata of an Instance which is one component of a
// Database deployment.
type Instance struct {
	ID        string         `json:"id"`
	Hostname  string         `json:"hostname"`
	Port      int            `json:"port"`
	Status    InstanceStatus `json:"status"`
	Type      InstanceType   `json:"type"`
	Features  []string       `json:"features"`
	PrivateIP string         `json:"private_ip"`
}

// DatabaseRes is the returned response from DatabaseShow
type DatabaseRes struct {
	Database Database `json:"database"`
}

// DatabaseShow returns the Database info of the given app/addonID
func (c *Client) DatabaseShow(ctx context.Context, app, addonID string) (Database, error) {
	var res DatabaseRes
	err := c.DBAPI(app, addonID).ResourceGet(ctx, "databases", addonID, nil, &res)
	if err != nil {
		return Database{}, errgo.Notef(err, "fail to get the database")
	}
	return res.Database, nil
}

// DatabaseUpdatePeriodicBackupsConfigParams contains the parameters which can
// be tweaked to update how periodic backups are triggered.
type DatabaseUpdatePeriodicBackupsConfigParams struct {
	ScheduledAt *int  `json:"periodic_backups_scheduled_at,omitempty"`
	Enabled     *bool `json:"periodic_backups_enabled,omitempty"`
}

// DatabaseUpdatePeriodicBackupsConfig updates the configuration of periodic
// backups for a given database addon
func (c *Client) DatabaseUpdatePeriodicBackupsConfig(ctx context.Context, app, addonID string, params DatabaseUpdatePeriodicBackupsConfigParams) (Database, error) {
	var dbRes DatabaseRes
	err := c.DBAPI(app, addonID).ResourceUpdate(ctx, "databases", addonID, map[string]DatabaseUpdatePeriodicBackupsConfigParams{
		"database": params,
	}, &dbRes)
	if err != nil {
		return Database{}, errgo.Notef(err, "fail to update periodic backups settings")
	}
	return dbRes.Database, nil
}

// DatabaseEnableFeatureParams contains the feature which has to be enabled
type DatabaseEnableFeatureParams struct {
	Feature DatabaseFeature `json:"feature"`
}

// DatabaseFeatureStatus is a type of string representing the advancement of
// the application of a database feature
type DatabaseFeatureStatus string

const (
	// DatabaseFeatureStatusActivated is set when the feature has been enabled with success
	DatabaseFeatureStatusActivated DatabaseFeatureStatus = "ACTIVATED"
	// DatabaseFeatureStatusPending is set when the feature is being enabled
	DatabaseFeatureStatusPending DatabaseFeatureStatus = "PENDING"
	// DatabaseFeatureStatusFailed is set when the feature failed to get enabeld
	DatabaseFeatureStatusFailed DatabaseFeatureStatus = "FAILED"
)

// DatabaseEnableFeatureResponse is the response structure from DatabaseEnableFeature
type DatabaseEnableFeatureResponse struct {
	Name    string                `json:"name"`
	Status  DatabaseFeatureStatus `json:"status"`
	Message string                `json:"message"`
}

// DatabaseEnableFeature enable a feature on a given database addon.
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

// DatabaseDisableFeatureResponse is the response body of DatabaseDisableFeature
type DatabaseDisableFeatureResponse struct {
	Message string `json:"message"`
}

// DatabaseDisableFeature disables a feature on a given database addon
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
