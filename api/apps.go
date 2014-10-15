package api

import (
	"net/http"
)

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

func AppsDestroy(id string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "DELETE",
		"endpoint": "/apps/" + id,
		"expected": Statuses{204},
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
		"expected": Statuses{201},
	}
	return Do(req)
}
