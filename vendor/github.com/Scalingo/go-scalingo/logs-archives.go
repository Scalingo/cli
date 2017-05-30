package scalingo

import (
	"encoding/json"
	"io/ioutil"

	errgo "gopkg.in/errgo.v1"
)

type LogsItem struct {
	Url  string `json:"url"`
	From string `json:"from"`
	To   string `json:"to"`
	Size int64  `json:"size"`
}

type LogsResponse struct {
	NextCursor string     `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
	Archives   []LogsItem `json:"archives"`
}

func (c *Client) LogsArchives(app string, cursor string) (*LogsResponse, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/logs_archives",
		Params: map[string]interface{}{
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

	var logsRes = LogsResponse{}
	err = json.Unmarshal(body, &logsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &logsRes, nil
}
