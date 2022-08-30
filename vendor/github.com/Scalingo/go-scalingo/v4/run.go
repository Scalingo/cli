package scalingo

import (
	"encoding/json"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v4/http"
)

type RunsService interface {
	Run(opts RunOpts) (*RunRes, error)
}

var _ RunsService = (*Client)(nil)

type RunOpts struct {
	App        string
	Command    []string
	Env        map[string]string
	Size       string
	Detached   bool
	HasUploads bool
}

type RunRes struct {
	Container *Container `json:"container"`
	AttachURL string     `json:"attach_url"`
}

func (c *Client) Run(opts RunOpts) (*RunRes, error) {
	req := &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + opts.App + "/run",
		Params: map[string]interface{}{
			"command":     strings.Join(opts.Command, " "),
			"env":         opts.Env,
			"size":        opts.Size,
			"detached":    opts.Detached,
			"has_uploads": opts.HasUploads,
		},
	}
	res, err := c.ScalingoAPI().Do(req)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var runRes RunRes
	err = json.NewDecoder(res.Body).Decode(&runRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &runRes, nil
}
