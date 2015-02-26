package api

import (
	"net/http"
	"strings"
)

func Run(app string, command []string, env map[string]string) (*http.Response, error) {
	req := &APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/run",
		Params: map[string]interface{}{
			"command": strings.Join(command, " "),
			"env":     env,
		},
	}
	return req.Do()
}
