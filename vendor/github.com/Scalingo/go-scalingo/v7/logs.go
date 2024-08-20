package scalingo

import (
	"context"
	"net/http"
	"net/url"

	"gopkg.in/errgo.v1"

	httpclient "github.com/Scalingo/go-scalingo/v7/http"
)

type LogsService interface {
	LogsURL(ctx context.Context, app string) (*http.Response, error)
	Logs(ctx context.Context, logsURL string, n int, filter string) (*http.Response, error)
}

var _ LogsService = (*Client)(nil)

func (c *Client) LogsURL(ctx context.Context, app string) (*http.Response, error) {
	req := &httpclient.APIRequest{
		Endpoint: "/apps/" + app + "/logs",
	}
	return c.ScalingoAPI().Do(ctx, req)
}

func (c *Client) Logs(ctx context.Context, logsURL string, n int, filter string) (*http.Response, error) {
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
			"token":     u.Query().Get("token"),
			"timestamp": u.Query().Get("timestamp"),
			"n":         n,
			"filter":    filter,
		},
	}
	return c.ScalingoAPI().Do(ctx, req)
}
