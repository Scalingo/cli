package scalingo

import (
	"time"

	errgo "gopkg.in/errgo.v1"
)

type DatabaseStatus string

const (
	DatabaseStatusRunning   DatabaseStatus = "running"
	DatabaseStatusMigrating DatabaseStatus = "migrating"
	DatabaseStatusUpgrading DatabaseStatus = "upgrading"
	DatabaseStatusStopped   DatabaseStatus = "stopped"
)

type Database struct {
	ID                         string              `json:"id"`
	CreatedAt                  time.Time           `json:"created_at"`
	ResourceID                 string              `json:"resource_id"`
	AppName                    string              `json:"app_name"`
	EncryptionAtRest           bool                `json:"encryption_at_rest"`
	Features                   []map[string]string `json:"features"`
	Plan                       string              `json:"plan"`
	Status                     DatabaseStatus      `json:"status"`
	TypeID                     string              `json:"type_id"`
	TypeName                   string              `json:"type_name"`
	VersionID                  string              `json:"version_id"`
	MongoReplSetName           string              `json:"mongo_repl_set_name"`
	Instances                  []Instance          `json:"instances"`
	NextVersionID              string              `json:"next_version_id"`
	ReadableVersion            string              `json:"readable_version"`
	Hostname                   string              `json:"hostname"`
	CurrentOperationID         string              `json:"current_operation_id"`
	Cluster                    bool                `json:"cluster"`
	PeriodicBackupsEnabled     bool                `json:"periodic_backups_enabled"`
	PeriodicBackupsScheduledAt int                 `json:"periodic_backups_scheduled_at"` // Hour of the day of the periodic backups (UTC)
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
	ID       string         `json:"id"`
	Hostname string         `json:"hostname"`
	Port     int            `json:"port"`
	Status   InstanceStatus `json:"status"`
	Type     InstanceType   `json:"type"`
	Features []string       `json:"features"`
	SandIP   string         `json:"sand_ip"`
}

type DatabaseRes struct {
	Database Database `json:"database"`
}

func (c *Client) DatabaseShow(app, addonID string) (Database, error) {
	var db Database
	err := c.DBAPI(app, addonID).ResourceGet("databases", addonID, nil, &db)
	if err != nil {
		return Database{}, errgo.Notef(err, "fail to get the database")
	}
	return db, nil
}

type PeriodicBackupsConfigParams struct {
	ScheduledAt *int  `json:"periodic_backups_scheduled_at,omitempty"`
	Enabled     *bool `json:"periodic_backups_enabled,omitempty"`
}

func (c *Client) PeriodicBackupsConfig(app, addonID string, params PeriodicBackupsConfigParams) (Database, error) {
	var dbRes DatabaseRes
	err := c.DBAPI(app, addonID).ResourceUpdate("databases", addonID, map[string]PeriodicBackupsConfigParams{
		"database": params,
	}, &dbRes)
	if err != nil {
		return Database{}, errgo.Notef(err, "fail to update periodic backups settings")
	}
	return dbRes.Database, nil
}
