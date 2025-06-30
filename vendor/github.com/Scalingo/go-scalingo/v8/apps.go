package scalingo

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v8/http"
)

type AppStatus string

const (
	AppStatusNew        AppStatus = "new"
	AppStatusRunning    AppStatus = "running"
	AppStatusStopped    AppStatus = "stopped"
	AppStatusScaling    AppStatus = "scaling"
	AppStatusRestarting AppStatus = "restarting"
)

type AppsService interface {
	AppsList(ctx context.Context) ([]*App, error)
	AppsShow(ctx context.Context, appName string) (*App, error)
	AppsDestroy(ctx context.Context, name string, currentName string) error
	AppsRename(ctx context.Context, name string, newName string) (*App, error)
	AppsTransfer(ctx context.Context, name string, email string) (*App, error)
	AppsSetStack(ctx context.Context, name string, stackID string) (*App, error)
	AppsRestart(ctx context.Context, app string, scope *AppsRestartParams) (*http.Response, error)
	AppsCreate(ctx context.Context, opts AppsCreateOpts) (*App, error)
	AppsStats(ctx context.Context, app string) (*AppStatsRes, error)
	AppsContainerTypes(ctx context.Context, app string) ([]ContainerType, error)
	AppsContainersPs(ctx context.Context, app string) ([]Container, error)
	AppsScale(ctx context.Context, app string, params *AppsScaleParams) (*http.Response, error)
	AppsForceHTTPS(ctx context.Context, name string, enable bool) (*App, error)
	AppsStickySession(ctx context.Context, name string, enable bool) (*App, error)
	AppsRouterLogs(ctx context.Context, name string, enable bool) (*App, error)
}

var _ AppsService = (*Client)(nil)

type ContainerType struct {
	AppID   string `json:"app_id"`
	Name    string `json:"name"`
	Amount  int    `json:"amount"`
	Command string `json:"command"`
	Size    string `json:"size"`
}

type ContainerStat struct {
	ID                 string `json:"id"`
	CPUUsage           int    `json:"cpu_usage"`
	MemoryUsage        int64  `json:"memory_usage"`
	SwapUsage          int64  `json:"swap_usage"`
	MemoryLimit        int64  `json:"memory_limit"`
	SwapLimit          int64  `json:"swap_limit"`
	HighestMemoryUsage int64  `json:"highest_memory_usage"`
	HighestSwapUsage   int64  `json:"highest_swap_usage"`
}

type AppStatsRes struct {
	Stats []*ContainerStat `json:"stats"`
}

type AppsScaleParams struct {
	Containers []ContainerType `json:"containers"`
}

type AppsPsRes struct {
	Containers []Container `json:"containers"`
}

type AppsContainerTypesRes struct {
	Containers []ContainerType `json:"containers"`
}

type AppsCreateOpts struct {
	Name      string `json:"name"`
	ParentApp string `json:"parent_id,omitempty"`
	StackID   string `json:"stack_id,omitempty"`
}

type AppResponse struct {
	App *App `json:"app"`
}

type AppsRestartParams struct {
	Scope []string `json:"scope"`
}

type AppLinks struct {
	DeploymentsStream string `json:"deployments_stream"`
}

type App struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Region             string                 `json:"region"`
	Owner              Owner                  `json:"owner"`
	GitURL             string                 `json:"git_url"`
	URL                string                 `json:"url"`
	BaseURL            string                 `json:"base_url"`
	Status             AppStatus              `json:"status"`
	LastDeployedAt     *time.Time             `json:"last_deployed_at"`
	LastDeployedBy     string                 `json:"last_deployed_by"`
	CreatedAt          *time.Time             `json:"created_at"`
	UpdatedAt          *time.Time             `json:"updated_at"`
	Links              *AppLinks              `json:"links"`
	StackID            string                 `json:"stack_id"`
	StickySession      bool                   `json:"sticky_session"`
	ForceHTTPS         bool                   `json:"force_https"`
	RouterLogs         bool                   `json:"router_logs"`
	DataAccessConsent  *DataAccessConsent     `json:"data_access_consent,omitempty"`
	Flags              map[string]bool        `json:"flags"`
	Limits             map[string]interface{} `json:"limits"`
	HDSResource        bool                   `json:"hds_resource"`
	PrivateNetworksIDs []string               `json:"private_networks_ids"`
}

func (app App) String() string {
	return app.Name
}

func (c *Client) AppsList(ctx context.Context) ([]*App, error) {
	appsMap := map[string][]*App{}

	req := &httpclient.APIRequest{
		Endpoint: "/apps",
	}

	err := c.ScalingoAPI().DoRequest(ctx, req, &appsMap)
	if err != nil {
		return []*App{}, errgo.Mask(err, errgo.Any)
	}
	return appsMap["apps"], nil
}

func (c *Client) AppsShow(ctx context.Context, appName string) (*App, error) {
	var appMap map[string]*App
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + appName,
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &appMap)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return appMap["app"], nil
}

func (c *Client) AppsDestroy(ctx context.Context, name string, currentName string) error {
	req := &httpclient.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + name,
		Expected: httpclient.Statuses{204},
		Params: map[string]interface{}{
			"current_name": currentName,
		},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AppsRename(ctx context.Context, name string, newName string) (*App, error) {
	var appRes *AppResponse
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + name + "/rename",
		Expected: httpclient.Statuses{200},
		Params: map[string]interface{}{
			"current_name": name,
			"new_name":     newName,
		},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &appRes)
	if err != nil {
		return nil, err
	}

	return appRes.App, nil
}

func (c *Client) AppsTransfer(ctx context.Context, name string, email string) (*App, error) {
	var appRes *AppResponse
	req := &httpclient.APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + name,
		Expected: httpclient.Statuses{200},
		Params: map[string]interface{}{
			"app": map[string]string{
				"owner": email,
			},
		},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &appRes)
	if err != nil {
		return nil, err
	}

	return appRes.App, nil
}

func (c *Client) AppsSetStack(ctx context.Context, app string, stackID string) (*App, error) {
	req := &httpclient.APIRequest{
		Method:   "PATCH",
		Endpoint: "/apps/" + app,
		Expected: httpclient.Statuses{200},
		Params: map[string]interface{}{
			"app": map[string]string{
				"stack_id": stackID,
			},
		},
	}

	var appRes AppResponse
	err := c.ScalingoAPI().DoRequest(ctx, req, &appRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to request Scalingo API")
	}

	return appRes.App, nil
}

func (c *Client) AppsRestart(ctx context.Context, app string, scope *AppsRestartParams) (*http.Response, error) {
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/restart",
		Expected: httpclient.Statuses{202},
		Params:   scope,
	}

	return c.ScalingoAPI().Do(ctx, req)
}

func (c *Client) AppsCreate(ctx context.Context, opts AppsCreateOpts) (*App, error) {
	var appRes *AppResponse
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps",
		Expected: httpclient.Statuses{201},
		Params:   map[string]interface{}{"app": opts},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &appRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return appRes.App, nil
}

func (c *Client) AppsStats(ctx context.Context, app string) (*AppStatsRes, error) {
	var stats AppStatsRes
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/stats",
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &stats)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &stats, nil
}

func (c *Client) AppsContainersPs(ctx context.Context, app string) ([]Container, error) {
	var containersRes AppsPsRes
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/ps",
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &containersRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to execute the GET request to list containers")
	}

	return containersRes.Containers, nil
}

func (c *Client) AppsContainerTypes(ctx context.Context, app string) ([]ContainerType, error) {
	var containerTypesRes AppsContainerTypesRes
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/containers",
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &containerTypesRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to execute the GET request to list container types")
	}

	return containerTypesRes.Containers, nil
}

func (c *Client) AppsScale(ctx context.Context, app string, params *AppsScaleParams) (*http.Response, error) {
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scale",
		Params:   params,
		// Return 200 if app is scaled before deployment
		// Otherwise async job is triggered, it's 202
		Expected: httpclient.Statuses{200, 202},
	}
	return c.ScalingoAPI().Do(ctx, req)
}

func (c *Client) AppsForceHTTPS(ctx context.Context, name string, enable bool) (*App, error) {
	return c.appsUpdate(ctx, name, map[string]interface{}{
		"force_https": enable,
	})
}

func (c *Client) AppsRouterLogs(ctx context.Context, name string, enable bool) (*App, error) {
	return c.appsUpdate(ctx, name, map[string]interface{}{
		"router_logs": enable,
	})
}

func (c *Client) AppsStickySession(ctx context.Context, name string, enable bool) (*App, error) {
	return c.appsUpdate(ctx, name, map[string]interface{}{
		"sticky_session": enable,
	})
}

func (c *Client) appsUpdate(ctx context.Context, name string, params map[string]interface{}) (*App, error) {
	var appRes *AppResponse
	req := &httpclient.APIRequest{
		Method:   "PUT",
		Endpoint: "/apps/" + name,
		Expected: httpclient.Statuses{200},
		Params:   params,
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &appRes)
	if err != nil {
		return nil, err
	}

	return appRes.App, nil
}
