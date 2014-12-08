package api

import (
	"net/http"
)

func LogsURL(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/apps/" + app + "/logs",
		"expected": Statuses{200},
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
		"expected": Statuses{200},
	}
	return Do(req)
}
