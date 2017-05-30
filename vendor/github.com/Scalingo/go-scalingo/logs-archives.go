package scalingo

import (
	"net/http"
)

func (c *Client) LogsArchives(app string, cursor string) (*http.Response, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]interface{}{
			"cursor": cursor,
		},
	}
	return req.Do()
}
