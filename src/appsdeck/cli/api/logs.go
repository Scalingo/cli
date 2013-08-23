package api

import (
	"appsdeck/cli/config"
	"net/http"
)

func Logs(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"host":     config.C["APPSDECK_LOG"],
		"endpoint": "/apps/" + app + "/logs",
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
