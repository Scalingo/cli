package scalingo

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v6/http"
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
	req := &http.APIRequest{
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]string{
			"cursor": cursor,
		},
	}

	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	var logsRes = LogsArchivesResponse{}
	err = json.Unmarshal(body, &logsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &logsRes, nil
}

func (c *Client) LogsArchives(ctx context.Context, app string, page int) (*LogsArchivesResponse, error) {
	if page < 1 {
		return nil, errgo.New("Page must be greater than 0.")
	}

	req := &http.APIRequest{
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]string{
			"page": strconv.FormatInt(int64(page), 10),
		},
	}

	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	var logsRes = LogsArchivesResponse{}
	err = json.Unmarshal(body, &logsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &logsRes, nil
}
