package scalingo

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	errgo "gopkg.in/errgo.v1"
)

type LogsArchivesService interface {
	LogsArchivesByCursor(app string, cursor string) (*LogsArchivesResponse, error)
	LogsArchives(app string, page int) (*LogsArchivesResponse, error)
}

var _ LogsArchivesService = (*Client)(nil)

type LogsArchiveItem struct {
	Url  string `json:"url"`
	From string `json:"from"`
	To   string `json:"to"`
	Size int64  `json:"size"`
}

type LogsArchivesResponse struct {
	NextCursor string            `json:"next_cursor"`
	HasMore    bool              `json:"has_more"`
	Archives   []LogsArchiveItem `json:"archives"`
}

func (c *Client) LogsArchivesByCursor(app string, cursor string) (*LogsArchivesResponse, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]string{
			"cursor": cursor,
		},
	}

	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	body, err := ioutil.ReadAll(res.Body)
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

func (c *Client) LogsArchives(app string, page int) (*LogsArchivesResponse, error) {
	if page < 1 {
		return nil, errgo.New("Page must be greater than 0.")
	}

	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]string{
			"page": strconv.FormatInt(int64(page), 10),
		},
	}

	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	body, err := ioutil.ReadAll(res.Body)
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
