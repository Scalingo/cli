package api

import (
	"net/http"
)

func Logs(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"host":     "10.1.0.2:10004",
		"endpoint": "/apps/" + app + "/logs",
	}
	return Do(req)
}

func LogsStream(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"host":     "10.1.0.2:10004",
		"endpoint": "/apps/" + app + "/logs/stream",
	}
	return Do(req)
}
