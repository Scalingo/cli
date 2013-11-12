package api

import (
	"appsdeck/config"
	"net/http"
)

func Logs(app string, n string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"host":     config.C["APPSDECK_LOG"],
		"endpoint": "/apps/" + app + "/logs",
		"params": map[string]interface{}{
			"n": n,
		},
	}
	return Do(req)
}

func LogsStream(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"host":     config.C["APPSDECK_LOG"],
		"endpoint": "/apps/" + app + "/logs/stream",
	}
	return Do(req)
}
