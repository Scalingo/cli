package scalingo

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"
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
		return nil, errors.Wrap(ctx, err, "create source")
	}

	return sourceRes.Source, nil
}
