package scalingo

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Scalingo/go-scalingo/v6/http"
	errors "github.com/Scalingo/go-utils/errors/v2"
)

type RunsService interface {
	Run(ctx context.Context, opts RunOpts) (*RunRes, error)
}

var _ RunsService = (*Client)(nil)

type RunOpts struct {
	App        string
	Command    []string
	Env        map[string]string
	Size       string
	Detached   bool
	Async      bool
	HasUploads bool
}

type RunRes struct {
	Container    *Container `json:"container"`
	AttachURL    string     `json:"attach_url"`
	OperationURL string     `json:"-"`
}

func (c *Client) Run(ctx context.Context, opts RunOpts) (*RunRes, error) {
	req := &http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + opts.App + "/run",
		Params: map[string]interface{}{
			"command":     strings.Join(opts.Command, " "),
			"env":         opts.Env,
			"size":        opts.Size,
			"detached":    opts.Detached,
			"async":       opts.Async,
			"has_uploads": opts.HasUploads,
		},
	}
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, errors.Notef(ctx, err, "request endpoint %v", req.Endpoint)
	}
	defer res.Body.Close()

	var runRes RunRes
	err = json.NewDecoder(res.Body).Decode(&runRes)
	if err != nil {
		return nil, errors.Notef(ctx, err, "decode response body")
	}

	runRes.OperationURL = res.Header.Get("Location")

	return &runRes, nil
}
