package api

import (
	"net/http"
)

func LogsURL(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/api/apps/" + app + "/logs",
	}
	return Do(req)
}

func Logs(url string, stream bool, n int) (*http.Response, error) {
	req := map[string]interface{}{
		"auth":   false,
		"method": "GET",
		"url":    url,
		"params": map[string]interface{}{
			"stream": stream,
			"n":      n,
		},
	}
	return Do(req)
}
