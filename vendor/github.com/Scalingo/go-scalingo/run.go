package scalingo

import (
	"net/http"
	"strings"
)

func (c *Client) Run(app string, command []string, env map[string]string, size string) (*http.Response, error) {
	req := &APIRequest{
		Client:   c,
		Method:   "POST",
		Endpoint: "/apps/" + app + "/run",
		Params: map[string]interface{}{
			"command": strings.Join(command, " "),
			"env":     env,
			"size":    size,
		},
	}
	return req.Do()
}
