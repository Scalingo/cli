package api

import (
	"net/http"
	"time"

	"gopkg.in/errgo.v1"
)

type Container struct {
	Name    string `json:"name"`
	Amount  int    `json:"amount"`
	Command string `json:"command"`
}

type AppsScaleParams struct {
	Containers []Container `json:"containers"`
}

type AppsPsRes struct {
	Containers []Container `json:"containers"`
}

type AppsRestartParams struct {
	Scope []string `json:"scope"`
}

type App struct {
	Id    string `json:"_id"`
	Name  string `json:"name"`
	Owner struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"owner"`
	GitUrl    string    `json:"git_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	LogsURL   string    `json:"logs_url"`
}

func (app App) String() string {
	return app.Name
}

type CreateAppParams struct {
	App *App `json:"app"`
}

func AppsList() (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps",
		"expected": Statuses{200},
	}
	return Do(req)
}

func AppsShow(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app,
		"expected": Statuses{200},
	}
	return Do(req)
}

func AppsDestroy(name string, currentName string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "DELETE",
		"endpoint": "/apps/" + name,
		"expected": Statuses{204},
		"params": map[string]interface{}{
			"current_name": currentName,
		},
	}
	return Do(req)
}

func AppsRestart(app string, scope *AppsRestartParams) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/restart",
		"expected": Statuses{202},
		"params":   scope,
	}
	return Do(req)
}

func AppsCreate(app string) (*App, int, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps",
		"params": map[string]interface{}{
			"app": map[string]interface{}{
				"name": app,
			},
		},
		"expected": Statuses{201},
	}
	res, err := Do(req)
	if err != nil {
		return nil, 0, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params *CreateAppParams
	err = ParseJSON(res, &params)
	if err != nil {
		return nil, res.StatusCode, errgo.Mask(err, errgo.Any)
	}

	return params.App, res.StatusCode, nil
}

func AppsPs(app string) ([]Container, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/containers",
		"expected": Statuses{200},
	}
	res, err := Do(req)
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

func AppsScale(app string, params *AppsScaleParams) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/scale",
		"params":   params,
		"expected": Statuses{202},
	}
	return Do(req)
}
