package api

import (
	"net/http"
	"strings"
)

func Run(app string, command []string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/api/apps/" + app + "/run",
		"params": map[string]interface{}{
			"command": strings.Join(command, " "),
		},
	}
	return Do(req)
}
