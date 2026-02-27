package scalingo

import (
	"context"
	"net/http"
	"net/url"

	httpclient "github.com/Scalingo/go-scalingo/v9/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type LogsService interface {
	LogsURL(ctx context.Context, app string) (*LogsURLRes, error)
	Logs(ctx context.Context, logsURL string, n int, filter string) (*http.Response, error)
}

var _ LogsService = (*Client)(nil)

type LogsURLRes struct {
	LogsURL string `json:"logs_url"`
	App     *App   `json:"app,omitempty"`
}

func (c *Client) LogsURL(ctx context.Context, app string) (*LogsURLRes, error) {
	var logsURLRes LogsURLRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "logs", nil, &logsURLRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get app logs URL")
	}

	return &logsURLRes, nil
}

func (c *Client) Logs(ctx context.Context, logsURL string, n int, filter string) (*http.Response, error) {
	u, err := url.Parse(logsURL)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parse logs URL")
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
