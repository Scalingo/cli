package api

import (
	"net/http"
	"strings"
)

func Run(app string, command []string, env map[string]string) (*http.Response, error) {
	req := map[string]interface{}{
		"method":   "POST",
		"endpoint": "/apps/" + app + "/run",
		"params": map[string]interface{}{
			"command": strings.Join(command, " "),
			"env":     env,
		},
		"expected": Statuses{200},
	}
	return Do(req)
}
