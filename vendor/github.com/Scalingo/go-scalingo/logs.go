package scalingo

import (
	"net/http"
	"net/url"

	"gopkg.in/errgo.v1"
)

type LogsService interface {
	LogsURL(app string) (*http.Response, error)
	Logs(logsURL string, n int, filter string) (*http.Response, error)
}

type LogsClient struct {
	*backendConfiguration
}

func (c *LogsClient) LogsURL(app string) (*http.Response, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
		Endpoint: "/apps/" + app + "/logs",
	}
	return req.Do()
}

func (c *LogsClient) Logs(logsURL string, n int, filter string) (*http.Response, error) {
	u, err := url.Parse(logsURL)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	req := &APIRequest{
		Client:   c.backendConfiguration,
		NoAuth:   true,
		Expected: Statuses{200, 204, 404},
		URL:      u.Scheme + "://" + u.Host,
		Endpoint: u.Path,
		Params: map[string]interface{}{
			"token":  u.Query().Get("token"),
			"n":      n,
			"filter": filter,
		},
	}
	return req.Do()
}
