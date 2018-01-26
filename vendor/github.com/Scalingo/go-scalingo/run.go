package scalingo

import (
	"strings"

	"gopkg.in/errgo.v1"
)

type RunsService interface {
	Run(opts RunOpts) (*RunRes, error)
}

type RunsClient struct {
	*backendConfiguration
}

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

func (c *RunsClient) Run(opts RunOpts) (*RunRes, error) {
	req := &APIRequest{
		Client:   c.backendConfiguration,
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
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var runRes RunRes
	err = ParseJSON(res, &runRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &runRes, nil
}
