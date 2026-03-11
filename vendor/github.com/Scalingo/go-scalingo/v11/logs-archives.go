package scalingo

import (
	"context"
	"strconv"

	"github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type LogsArchivesService interface {
	LogsArchivesByCursor(ctx context.Context, app string, cursor string) (*LogsArchivesResponse, error)
	LogsArchives(ctx context.Context, app string, page int) (*LogsArchivesResponse, error)
}

var _ LogsArchivesService = (*Client)(nil)

type LogsArchiveItem struct {
	URL  string `json:"url"`
	From string `json:"from"`
	To   string `json:"to"`
	Size int64  `json:"size"`
}

type LogsArchivesResponse struct {
	NextCursor string            `json:"next_cursor"`
	HasMore    bool              `json:"has_more"`
	Archives   []LogsArchiveItem `json:"archives"`
}

func (c *Client) LogsArchivesByCursor(ctx context.Context, app string, cursor string) (*LogsArchivesResponse, error) {
	var logsRes LogsArchivesResponse
	req := &http.APIRequest{
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]string{
			"cursor": cursor,
		},
	}

	err := c.ScalingoAPI().DoRequest(ctx, req, &logsRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list logs archives by cursor")
	}

	return &logsRes, nil
}

func (c *Client) LogsArchives(ctx context.Context, app string, page int) (*LogsArchivesResponse, error) {
	if page < 1 {
		return nil, errors.New(ctx, "Page must be greater than 0.")
	}

	req := &http.APIRequest{
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]string{
			"page": strconv.FormatInt(int64(page), 10),
		},
	}

	var logsRes LogsArchivesResponse
	err := c.ScalingoAPI().DoRequest(ctx, req, &logsRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list logs archives")
	}

	return &logsRes, nil
}
