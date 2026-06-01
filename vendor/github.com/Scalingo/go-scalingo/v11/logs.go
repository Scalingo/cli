package scalingo

import (
	"context"
	stderrors "errors"
	"io"
	"net/http"
	"net/url"

	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type LogsService interface {
	LogsURL(ctx context.Context, app string) (*LogsURLRes, error)
	// Logs returns the raw http.Response from the request to the API. This response body contains the requested log lines in raw text.
	// It has been decided to let the user of this function decides how to best read the body (type is io.ReadCloser) depending on their context.
	Logs(ctx context.Context, logsURL string, n int, filter string) (io.ReadCloser, error)
}

var _ LogsService = (*Client)(nil)

type LogsURLRes struct {
	LogsURL string `json:"logs_url"`
	App     *App   `json:"app,omitempty"`
}

func (c *Client) LogsURL(ctx context.Context, app string) (*LogsURLRes, error) {
	var logsURLRes LogsURLRes
	err := c.ScalingoAPI().SubresourceList(ctx, appsResource, app, logsResource, nil, &logsURLRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get app logs URL")
	}

	return &logsURLRes, nil
}

var ErrNoLogs = stderrors.New("application didn't logged anything yet")

func (c *Client) Logs(ctx context.Context, logsURL string, n int, filter string) (io.ReadCloser, error) {
	u, err := url.Parse(logsURL)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parse logs URL")
	}
	req := &httpclient.APIRequest{
		NoAuth:   true,
		Expected: httpclient.Statuses{http.StatusOK, http.StatusNoContent, http.StatusNotFound},
		URL:      u.Scheme + "://" + u.Host,
		Endpoint: u.Path,
		Params: map[string]any{
			"token":     u.Query().Get("token"),
			"timestamp": u.Query().Get("timestamp"),
			"n":         n,
			"filter":    filter,
		},
	}
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "request Scalingo API to get the application logs")
	}

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusNoContent {
		return nil, ErrNoLogs
	}

	return res.Body, nil
}
