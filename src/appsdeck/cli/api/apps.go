package api

import (
	"net/http"
)

func AppsList() (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "GET",
		"endpoint": "/api/apps",
	}

	res, err := Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
