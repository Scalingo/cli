package scalingo

import (
	"context"
	"net/http"
	"strings"

	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
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
	HasUploads bool
}

type RunRes struct {
	Container    *Container `json:"container"`
	AttachURL    string     `json:"attach_url"`
	OperationURL string     `json:"operation_url"`
}

func (c *Client) Run(ctx context.Context, opts RunOpts) (*RunRes, error) {
	var runRes RunRes
	req := &httpclient.APIRequest{
		Method:   http.MethodPost,
		Endpoint: "/apps/" + opts.App + "/run",
		Params: map[string]any{
			"command":     strings.Join(opts.Command, " "),
			"env":         opts.Env,
			"size":        opts.Size,
			"detached":    opts.Detached,
			"has_uploads": opts.HasUploads,
		},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &runRes)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "request endpoint %v", req.Endpoint)
	}

	return &runRes, nil
}
