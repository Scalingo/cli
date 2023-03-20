package scalingo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v6/http"
)

type DeploymentsService interface {
	DeploymentList(ctx context.Context, app string) ([]*Deployment, error)
	DeploymentListWithPagination(ctx context.Context, app string, opts PaginationOpts) ([]*Deployment, PaginationMeta, error)
	Deployment(ctx context.Context, app string, deploy string) (*Deployment, error)
	DeploymentLogs(ctx context.Context, deployURL string) (*http.Response, error)
	DeploymentStream(ctx context.Context, deployURL string) (*websocket.Conn, error)
	DeploymentsCreate(ctx context.Context, app string, params *DeploymentsCreateParams) (*Deployment, error)
}

var _ DeploymentsService = (*Client)(nil)

type DeploymentStatus string

const (
	StatusSuccess      DeploymentStatus = "success"
	StatusQueued       DeploymentStatus = "queued"
	StatusBuilding     DeploymentStatus = "building"
	StatusStarting     DeploymentStatus = "starting"
	StatusPushing      DeploymentStatus = "pushing"
	StatusAborted      DeploymentStatus = "aborted"
	StatusBuildError   DeploymentStatus = "build-error"
	StatusCrashedError DeploymentStatus = "crashed-error"
	StatusTimeoutError DeploymentStatus = "timeout-error"
	StatusHookError    DeploymentStatus = "hook-error"
)

type Deployment struct {
	ID             string           `json:"id"`
	AppID          string           `json:"app_id"`
	CreatedAt      *time.Time       `json:"created_at"`
	Status         DeploymentStatus `json:"status"`
	GitRef         string           `json:"git_ref"`
	Image          string           `json:"image"`
	Registry       string           `json:"registry"`
	Duration       int              `json:"duration"`
	PostdeployHook string           `json:"postdeploy_hook"`
	ImageSize      uint64           `json:"image_size"`
	StackBaseImage string           `json:"stack_base_image"`
	User           *User            `json:"pusher"`
	Links          *DeploymentLinks `json:"links"`
}

// DeploymentEventType holds all different deployment stream types of event.
type DeploymentEventType string

const (
	DeploymentEventTypePing   DeploymentEventType = "ping"
	DeploymentEventTypeNew    DeploymentEventType = "new"
	DeploymentEventTypeLog    DeploymentEventType = "log"
	DeploymentEventTypeStatus DeploymentEventType = "status"
)

// DeploymentEvent represents a deployment stream event sent on the websocket.
type DeploymentEvent struct {
	// ID of the deployment which this event belongs to
	ID   string              `json:"id"`
	Type DeploymentEventType `json:"type"`
	Data json.RawMessage     `json:"data"`
}

// DeploymentEventDataLog is the data type present in the DeploymentEvent.Data field if the DeploymentEvent.Type is DeploymentEventDataLog
type DeploymentEventDataLog struct {
	Content string `json:"content"`
}

// DeploymentEventDataStatus is the data type present in the DeploymentEvent.Data field if the DeploymentEvent.Type is DeploymentEventDataStatus
type DeploymentEventDataStatus struct {
	Status DeploymentStatus `json:"status"`
}

type DeploymentsCreateParams struct {
	GitRef    *string `json:"git_ref"`
	SourceURL string  `json:"source_url"`
}

type DeploymentsCreateRes struct {
	Deployment *Deployment `json:"deployment"`
}

func (d *Deployment) HasFailed() bool {
	return HasFailedString(d.Status)
}

func HasFailedString(status DeploymentStatus) bool {
	if !IsFinishedString(status) {
		return false
	}

	if status == StatusSuccess {
		return false
	}
	return true
}

func (d *Deployment) IsFinished() bool {
	return IsFinishedString(d.Status)
}

func IsFinishedString(status DeploymentStatus) bool {
	return status != StatusBuilding && status != StatusStarting &&
		status != StatusPushing && status != StatusQueued
}

type DeploymentList struct {
	Deployments []*Deployment `json:"deployments"`
	Meta        struct {
		PaginationMeta PaginationMeta `json:"pagination"`
	}
}

type DeploymentLinks struct {
	Output string `json:"output"`
}

type AuthenticationData struct {
	Token string `json:"token"`
}

type AuthStruct struct {
	Type string             `json:"type"`
	Data AuthenticationData `json:"data"`
}

func (c *Client) DeploymentList(ctx context.Context, app string) ([]*Deployment, error) {
	var deployments DeploymentList
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/deployments",
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &deployments)
	if err != nil {
		return []*Deployment{}, errgo.Notef(err, "fail to list the deployments")
	}

	return deployments.Deployments, nil
}

func (c *Client) DeploymentListWithPagination(ctx context.Context, app string, opts PaginationOpts) ([]*Deployment, PaginationMeta, error) {
	var deployments DeploymentList
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "deployments", opts.ToMap(), &deployments)
	if err != nil {
		return []*Deployment{}, PaginationMeta{}, errgo.Notef(err, "fail to list the deployments with pagination")
	}

	return deployments.Deployments, deployments.Meta.PaginationMeta, nil
}

func (c *Client) Deployment(ctx context.Context, app string, deploy string) (*Deployment, error) {
	var deploymentMap map[string]*Deployment
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/deployments/" + deploy,
	}

	err := c.ScalingoAPI().DoRequest(ctx, req, &deploymentMap)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return deploymentMap["deployment"], nil
}

func (c *Client) DeploymentLogs(ctx context.Context, deployURL string) (*http.Response, error) {
	u, err := url.Parse(deployURL)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	req := &httpclient.APIRequest{
		Expected: httpclient.Statuses{200, 404},
		Endpoint: u.Path,
		URL:      u.Scheme + "://" + u.Host,
	}

	return c.ScalingoAPI().Do(ctx, req)
}

// DeploymentStream returns a websocket connection to follow the various deployment events happening on an application. The type of the data sent on this connection is DeployEvent.
func (c *Client) DeploymentStream(ctx context.Context, deployURL string) (*websocket.Conn, error) {
	token, err := c.ScalingoAPI().TokenGenerator().GetAccessToken(ctx)
	if err != nil {
		return nil, errgo.Notef(err, "fail to generate token")
	}

	header := http.Header{}
	header.Add("Origin", "http://scalingo-cli.local/1")
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, deployURL, header)
	if err != nil {
		return nil, errgo.Notef(err, "fail to dial on url %s", deployURL)
	}
	defer resp.Body.Close()

	err = conn.WriteJSON(&AuthStruct{
		Type: "auth",
		Data: AuthenticationData{
			Token: token,
		},
	})
	if err != nil {
		return nil, errgo.Notef(err, "fail to write JSON, there must be an authentication issue")
	}

	return conn, nil
}

func (c *Client) DeploymentsCreate(ctx context.Context, app string, params *DeploymentsCreateParams) (*Deployment, error) {
	var response *DeploymentsCreateRes
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/deployments",
		Expected: httpclient.Statuses{201},
		Params: map[string]interface{}{
			"deployment": params,
		},
	}

	err := c.ScalingoAPI().DoRequest(ctx, req, &response)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return response.Deployment, nil
}

func (c *Client) DeploymentCacheReset(ctx context.Context, app string) error {
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/caches/deployment",
		Method:   "DELETE",
		Expected: httpclient.Statuses{204},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, nil)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
