package api

import (
	"net/http"
)

func Logs(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"host":    "10.1.0.2:10004",
		"endpoint": "/apps/"+app+"/logs",
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func LogsStream(app string) (*http.Response, error) {
	req := map[string]interface{}{
		"method" : "GET",
		"host": "10.1.0.2:10004",
		"endpoint": "/apps/"+app+"/logs/stream",
	}
	res, err := Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
