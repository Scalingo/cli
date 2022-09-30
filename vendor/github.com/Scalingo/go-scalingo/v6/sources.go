package scalingo

import (
	"context"

	"gopkg.in/errgo.v1"
)

type SourcesService interface {
	SourcesCreate(context.Context) (*Source, error)
}

var _ SourcesService = (*Client)(nil)

type SourcesCreateResponse struct {
	Source *Source `json:"source"`
}

type Source struct {
	DownloadURL string `json:"download_url"`
	UploadURL   string `json:"upload_url"`
}

func (c *Client) SourcesCreate(ctx context.Context) (*Source, error) {
	var sourceRes SourcesCreateResponse
	err := c.ScalingoAPI().ResourceAdd(ctx, "sources", nil, &sourceRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return sourceRes.Source, nil
}
