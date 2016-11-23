package scalingo

import (
	"net/http"
	"time"

	"gopkg.in/errgo.v1"
)

type ContainerType struct {
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
	Containers []ContainerType `json:"containers"`
}

type AppsCreateOpts struct {
	Name      string `json:"name"`
	ParentApp string `json:"parent_id"`
}

type AppsRestartParams struct {
	Scope []string `json:"scope"`
}

type AppLinks struct {
	DeploymentsStream string `json:"deployments_stream"`
}

type AppDomains struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	SSL  bool   `json:"ssl"`
}

type App struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Owner struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Billable bool   `json:"billable"`
	} `json:"owner"`
	GitUrl         string        `json:"git_url"`
	LastDeployedAt *time.Time    `json:"last_deployed_at"`
	LastDeployedBy string        `json:"last_deployed_by"`
	CreatedAt      *time.Time    `json:"created_at"`
	UpdatedAt      *time.Time    `json:"update_at"`
	Links          *AppLinks     `json:"links"`
	Domains        []*AppDomains `json:"domains"`
}

func (app App) String() string {
	return app.Name
}

type CreateAppParams struct {
	App *App `json:"app"`
}

func (c *Client) AppsList() ([]*App, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps",
	}

	res, err := req.Do()
	if err != nil {
		return []*App{}, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	appsMap := map[string][]*App{}
	err = ParseJSON(res, &appsMap)
	if err != nil {
		return []*App{}, errgo.Mask(err, errgo.Any)
	}
	return appsMap["apps"], nil
}

func (c *Client) AppsShow(appName string) (*App, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + appName,
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()
	var appMap map[string]*App
	err = ParseJSON(res, &appMap)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return appMap["app"], nil
}

func (c *Client) AppsDestroy(name string, currentName string) error {
	req := &APIRequest{
		Client:   c,
		Method:   "DELETE",
		Endpoint: "/apps/" + name,
		Expected: Statuses{204},
		Params: map[string]interface{}{
			"current_name": currentName,
		},
	}
	res, err := req.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (c *Client) AppsRestart(app string, scope *AppsRestartParams) (*http.Response, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "POST",
		Endpoint: "/apps/" + app + "/restart",
		Expected: Statuses{202},
		Params:   scope,
	}
	return req.Do()
}

func (c *Client) AppsCreate(opts AppsCreateOpts) (*App, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "POST",
		Endpoint: "/apps",
		Expected: Statuses{201},
		Params:   map[string]interface{}{"app": opts},
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params *CreateAppParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return params.App, nil
}

func (c *Client) AppsStats(app string) (*AppStatsRes, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/stats",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	var stats AppStatsRes
	err = ParseJSON(res, &stats)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &stats, nil
}

func (c *Client) AppsPs(app string) ([]ContainerType, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/containers",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	var containersRes AppsPsRes
	err = ParseJSON(res, &containersRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return containersRes.Containers, nil
}

func (c *Client) AppsScale(app string, params *AppsScaleParams) (*http.Response, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scale",
		Params:   params,
		Expected: Statuses{202},
	}
	return req.Do()
}
