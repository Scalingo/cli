package api

import (
	"net/http"
)

func AppsList() (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps",
	}
	return Do(req)
}

func AppsShow(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app,
	}
	return Do(req)
}

func AppsDestroy(id string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "DELETE",
		"endpoint": "/apps/" + id,
	}
	return Do(req)
}

func AppsCreate(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps",
		"params": map[string]interface{}{
			"app": map[string]interface{}{
				"name": app,
			},
		},
	}
	return Do(req)
}
