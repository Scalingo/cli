package scalingo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	httpclient "github.com/Scalingo/go-scalingo/http"
	"golang.org/x/net/websocket"
	"gopkg.in/errgo.v1"
)

type DeploymentsService interface {
	DeploymentList(app string) ([]*Deployment, error)
	Deployment(app string, deploy string) (*Deployment, error)
	DeploymentLogs(deployURL string) (*http.Response, error)
	DeploymentStream(deployURL string) (*websocket.Conn, error)
	DeploymentsCreate(app string, params *DeploymentsCreateParams) (*Deployment, error)
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
	User           *User            `json:"pusher"`
	Links          *DeploymentLinks `json:"links"`
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

func (c *Client) DeploymentList(app string) ([]*Deployment, error) {
	var deployments DeploymentList
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/deployments",
	}
	err := c.ScalingoAPI().DoRequest(req, &deployments)
	if err != nil {
		return []*Deployment{}, errgo.Mask(err, errgo.Any)
	}

	return deployments.Deployments, nil
}

func (c *Client) Deployment(app string, deploy string) (*Deployment, error) {
	var deploymentMap map[string]*Deployment
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/deployments/" + deploy,
	}

	err := c.ScalingoAPI().DoRequest(req, &deploymentMap)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return deploymentMap["deployment"], nil
}

func (c *Client) DeploymentLogs(deployURL string) (*http.Response, error) {
	u, err := url.Parse(deployURL)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	req := &httpclient.APIRequest{
		Expected: httpclient.Statuses{200, 404},
		Endpoint: u.Path,
		URL:      u.Scheme + "://" + u.Host,
	}

	return c.ScalingoAPI().Do(req)
}

func (c *Client) DeploymentStream(deployURL string) (*websocket.Conn, error) {
	token, err := c.ScalingoAPI().TokenGenerator().GetAccessToken()
	if err != nil {
		return nil, errgo.Notef(err, "fail to generate token")
	}
	authString, err := json.Marshal(&AuthStruct{
		Type: "auth",
		Data: AuthenticationData{
			Token: token,
		},
	})
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	conn, err := websocket.Dial(deployURL, "", "http://scalingo-cli.local/1")
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	_, err = conn.Write(authString)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return conn, nil
}

func (c *Client) DeploymentsCreate(app string, params *DeploymentsCreateParams) (*Deployment, error) {
	var response *DeploymentsCreateRes
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/deployments",
		Expected: httpclient.Statuses{201},
		Params: map[string]interface{}{
			"deployment": params,
		},
	}

	err := c.ScalingoAPI().DoRequest(req, &response)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return response.Deployment, nil
}

func (c *Client) DeploymentCacheReset(app string) error {
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/caches/deployment",
		Method:   "DELETE",
		Expected: httpclient.Statuses{204},
	}
	err := c.ScalingoAPI().DoRequest(req, nil)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
