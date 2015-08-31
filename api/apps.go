package api

import (
	"net/http"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

type Container struct {
	Name    string `json:"name"`
	Amount  int    `json:"amount"`
	Command string `json:"command"`
	Size    string `json:"size"`
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
}

func (app App) String() string {
	return app.Name
}

type CreateAppParams struct {
	App *App `json:"app"`
}

func AppsList() ([]*App, error) {
	req := &APIRequest{
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

func AppsShow(app string) (*http.Response, error) {
	req := &APIRequest{
		Endpoint: "/apps/" + app,
	}
	return req.Do()
}

func AppsDestroy(name string, currentName string) (*http.Response, error) {
	req := &APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + name,
		Expected: Statuses{204},
		Params: map[string]interface{}{
			"current_name": currentName,
		},
	}
	return req.Do()
}

func AppsRestart(app string, scope *AppsRestartParams) (*http.Response, error) {
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/restart",
		Expected: Statuses{202},
		Params:   scope,
	}
	return req.Do()
}

func AppsCreate(app string) (*App, error) {
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/apps",
		Expected: Statuses{201},
		Params: map[string]interface{}{
			"app": map[string]interface{}{
				"name": app,
			},
		},
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

func AppsPs(app string) ([]Container, error) {
	req := &APIRequest{
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

// Handling 422 error as not a standard 422 error
// {
//    "errors": {
//       "web": {
//          "amount": ["is negative"]
//       }
//    }
// }
func AppsScale(app string, params *AppsScaleParams) (*http.Response, error) {
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scale",
		Params:   params,
		Expected: Statuses{202, 422},
	}
	return req.Do()
}
