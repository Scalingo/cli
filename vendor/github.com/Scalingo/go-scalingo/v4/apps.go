package scalingo

import (
	"net/http"
	"time"

	httpclient "github.com/Scalingo/go-scalingo/v4/http"
	"gopkg.in/errgo.v1"
)

type AppStatus string

const (
	AppStatusNew        = AppStatus("new")
	AppStatusRunning    = AppStatus("running")
	AppStatusStopped    = AppStatus("stopped")
	AppStatusScaling    = AppStatus("scaling")
	AppStatusRestarting = AppStatus("restarting")
)

type AppsService interface {
	AppsList() ([]*App, error)
	AppsShow(appName string) (*App, error)
	AppsDestroy(name string, currentName string) error
	AppsRename(name string, newName string) (*App, error)
	AppsTransfer(name string, email string) (*App, error)
	AppsSetStack(name string, stackID string) (*App, error)
	AppsRestart(app string, scope *AppsRestartParams) (*http.Response, error)
	AppsCreate(opts AppsCreateOpts) (*App, error)
	AppsStats(app string) (*AppStatsRes, error)
	AppsContainerTypes(app string) ([]ContainerType, error)
	// Deprecated: Use AppsContainerTypes instead
	AppsPs(app string) ([]ContainerType, error)
	AppsContainersPs(app string) ([]Container, error)
	AppsScale(app string, params *AppsScaleParams) (*http.Response, error)
	AppsForceHTTPS(name string, enable bool) (*App, error)
	AppsStickySession(name string, enable bool) (*App, error)
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
	CpuUsage           int    `json:"cpu_usage"`
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
	Id     string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Owner  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"owner"`
	GitUrl         string                 `json:"git_url"`
	Url            string                 `json:"url"`
	BaseURL        string                 `json:"base_url"`
	Status         AppStatus              `json:"status"`
	LastDeployedAt *time.Time             `json:"last_deployed_at"`
	LastDeployedBy string                 `json:"last_deployed_by"`
	CreatedAt      *time.Time             `json:"created_at"`
	UpdatedAt      *time.Time             `json:"updated_at"`
	Links          *AppLinks              `json:"links"`
	StackID        string                 `json:"stack_id"`
	StickySession  bool                   `json:"sticky_session"`
	ForceHTTPS     bool                   `json:"force_https"`
	RouterLogs     bool                   `json:"router_logs"`
	Flags          map[string]bool        `json:"flags"`
	Limits         map[string]interface{} `json:"limits"`
}

func (app App) String() string {
	return app.Name
}

func (c *Client) AppsList() ([]*App, error) {
	appsMap := map[string][]*App{}

	req := &httpclient.APIRequest{
		Endpoint: "/apps",
	}

	err := c.ScalingoAPI().DoRequest(req, &appsMap)
	if err != nil {
		return []*App{}, errgo.Mask(err, errgo.Any)
	}
	return appsMap["apps"], nil
}

func (c *Client) AppsShow(appName string) (*App, error) {
	var appMap map[string]*App
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + appName,
	}
	err := c.ScalingoAPI().DoRequest(req, &appMap)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return appMap["app"], nil
}

func (c *Client) AppsDestroy(name string, currentName string) error {
	req := &httpclient.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + name,
		Expected: httpclient.Statuses{204},
		Params: map[string]interface{}{
			"current_name": currentName,
		},
	}
	err := c.ScalingoAPI().DoRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AppsRename(name string, newName string) (*App, error) {
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
	err := c.ScalingoAPI().DoRequest(req, &appRes)
	if err != nil {
		return nil, err
	}

	return appRes.App, nil
}

func (c *Client) AppsTransfer(name string, email string) (*App, error) {
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
	err := c.ScalingoAPI().DoRequest(req, &appRes)
	if err != nil {
		return nil, err
	}

	return appRes.App, nil
}

func (c *Client) AppsSetStack(app string, stackID string) (*App, error) {
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
	err := c.ScalingoAPI().DoRequest(req, &appRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to request Scalingo API")
	}

	return appRes.App, nil
}

func (c *Client) AppsRestart(app string, scope *AppsRestartParams) (*http.Response, error) {
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/restart",
		Expected: httpclient.Statuses{202},
		Params:   scope,
	}

	return c.ScalingoAPI().Do(req)
}

func (c *Client) AppsCreate(opts AppsCreateOpts) (*App, error) {
	var appRes *AppResponse
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps",
		Expected: httpclient.Statuses{201},
		Params:   map[string]interface{}{"app": opts},
	}
	err := c.ScalingoAPI().DoRequest(req, &appRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return appRes.App, nil
}

func (c *Client) AppsStats(app string) (*AppStatsRes, error) {
	var stats AppStatsRes
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/stats",
	}
	err := c.ScalingoAPI().DoRequest(req, &stats)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &stats, nil
}

// Deprecated: Use AppsContainerTypes instead
func (c *Client) AppsPs(app string) ([]ContainerType, error) {
	return c.AppsContainerTypes(app)
}

func (c *Client) AppsContainersPs(app string) ([]Container, error) {
	var containersRes AppsPsRes
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/ps",
	}
	err := c.ScalingoAPI().DoRequest(req, &containersRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to execute the GET request to list containers")
	}

	return containersRes.Containers, nil
}

func (c *Client) AppsContainerTypes(app string) ([]ContainerType, error) {
	var containerTypesRes AppsContainerTypesRes
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/containers",
	}
	err := c.ScalingoAPI().DoRequest(req, &containerTypesRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to execute the GET request to list container types")
	}

	return containerTypesRes.Containers, nil
}

func (c *Client) AppsScale(app string, params *AppsScaleParams) (*http.Response, error) {
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scale",
		Params:   params,
		// Return 200 if app is scaled before deployment
		// Otherwise async job is triggered, it's 202
		Expected: httpclient.Statuses{200, 202},
	}
	return c.ScalingoAPI().Do(req)
}

func (c *Client) AppsForceHTTPS(name string, enable bool) (*App, error) {
	return c.appsUpdate(name, map[string]interface{}{
		"force_https": enable,
	})
}

func (c *Client) AppsStickySession(name string, enable bool) (*App, error) {
	return c.appsUpdate(name, map[string]interface{}{
		"sticky_session": enable,
	})
}

func (c *Client) appsUpdate(name string, params map[string]interface{}) (*App, error) {
	var appRes *AppResponse
	req := &httpclient.APIRequest{
		Method:   "PUT",
		Endpoint: "/apps/" + name,
		Expected: httpclient.Statuses{200},
		Params:   params,
	}
	err := c.ScalingoAPI().DoRequest(req, &appRes)
	if err != nil {
		return nil, err
	}

	return appRes.App, nil
}
