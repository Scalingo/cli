package api

import (
	"net/http"
)

func Run(app string, command []string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/run",
		"host": "10.1.0.2:10006",
	}
	return Do(req)
}
