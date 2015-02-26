package api

import (
	"net/http"
)

func LogsURL(app string) (*http.Response, error) {
	req := &APIRequest{
		Endpoint: "/apps/" + app + "/logs",
	}
	return req.Do()
}

func Logs(url string, stream bool, n int) (*http.Response, error) {
	req := &APIRequest{
		NoAuth: true,
		URL:    url,
		Params: map[string]interface{}{
			"stream": stream,
			"n":      n,
		},
	}
	return req.Do()
}
