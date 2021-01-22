package scalingo

import (
	"net/http"
	"net/url"

	httpclient "github.com/Scalingo/go-scalingo/v4/http"
	"gopkg.in/errgo.v1"
)

type LogsService interface {
	LogsURL(app string) (*http.Response, error)
	Logs(logsURL string, n int, filter string) (*http.Response, error)
}

var _ LogsService = (*Client)(nil)

func (c *Client) LogsURL(app string) (*http.Response, error) {
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/logs",
	}
	return c.ScalingoAPI().Do(req)
}

func (c *Client) Logs(logsURL string, n int, filter string) (*http.Response, error) {
	u, err := url.Parse(logsURL)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	req := &httpclient.APIRequest{
		NoAuth:   true,
		Expected: httpclient.Statuses{200, 204, 404},
		URL:      u.Scheme + "://" + u.Host,
		Endpoint: u.Path,
		Params: map[string]interface{}{
			"token":  u.Query().Get("token"),
			"n":      n,
			"filter": filter,
		},
	}
	return c.ScalingoAPI().Do(req)
}
